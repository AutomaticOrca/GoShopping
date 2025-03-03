package server

import (
	"database/sql"
	"errors"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/AutomaticOrca/GoShopping/internal/worker"
	util "github.com/AutomaticOrca/GoShopping/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type createOrderRequest struct {
	UserID int32 `json:"user_id" binding:"required"`
	Items  []struct {
		ProductID int32 `json:"product_id" binding:"required"`
		Quantity  int   `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required"`
}

type orderResponse struct {
	ID         int32     `json:"id"`
	UserID     int32     `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func newOrderResponse(order db.CreateOrderRow, userID pgtype.Int4, totalPrice pgtype.Numeric) orderResponse {
	return orderResponse{
		ID:         order.ID,
		UserID:     userID.Int32,                        // **pgtype.Int4 -> int32**
		TotalPrice: db.MustNumericToFloat64(totalPrice), // **pgtype.Numeric -> float64**
		CreatedAt:  order.CreatedAt.Time,                // **pgtype.Timestamp -> time.Time**
		Status:     "pending",
	}
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var totalPrice float64
	orderItems := []db.CreateOrderItemParams{}

	for _, item := range req.Items {
		product, err := server.store.GetProductByID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		convertedPrice := db.MustNumericToFloat64(product.Price)

		totalPrice += float64(item.Quantity) * convertedPrice

		orderItems = append(orderItems, db.CreateOrderItemParams{
			OrderID:   db.IntToInt4(0),
			ProductID: db.Int32ToInt4(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     db.Float64ToNumeric(convertedPrice),
			Subtotal:  db.Float64ToNumeric(float64(item.Quantity) * convertedPrice),
		})
	}

	totalPricePG := db.Float64ToNumeric(totalPrice)

	userIDPG := db.Int32ToInt4(req.UserID)

	arg := db.CreateOrderParams{
		UserID:     userIDPG,
		TotalPrice: totalPricePG,
	}

	order, err := server.store.CreateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for i := range orderItems {
		orderItems[i].OrderID = db.Int32ToInt4(order.ID)

		err = server.store.CreateOrderItem(ctx, orderItems[i])
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	taskPayload := &worker.PayloadOrderCancel{
		OrderID: order.ID,
	}

	err = server.taskDistributor.DistributeTaskOrderCancel(ctx, taskPayload, worker.OrderCancelDelay)
	if err != nil {
		log.Errorf("failed to enqueue order cancel task: %v", err)
	}

	rsp := newOrderResponse(order, userIDPG, totalPricePG)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getOrder(ctx *gin.Context) {
	orderID, err := util.ParseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := server.store.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orderItems, err := server.store.GetOrderItemsByOrderID(ctx, db.Int32ToInt4(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userID, err := util.GetUserIDFromContext(ctx) // 从 JWT 获取用户 ID
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if order.UserID.Int32 != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to access this order"})
		return
	}

	rsp := gin.H{
		"order": gin.H{
			"id":          order.ID,
			"user_id":     order.UserID.Int32,
			"total_price": db.MustNumericToFloat64(order.TotalPrice),
			"status":      order.Status,
			"created_at":  order.CreatedAt.Time,
		},
		"items": orderItems,
	}

	ctx.JSON(http.StatusOK, rsp)
}
