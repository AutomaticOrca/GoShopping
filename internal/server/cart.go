package server

import (
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/AutomaticOrca/GoShopping/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cart item response structure
type cartItemResponse struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	AddedAt     string  `json:"added_at"`
}

// Convert db.GetCartRow to cartItemResponse
func newCartItemResponse(item db.GetCartRow) cartItemResponse {
	return cartItemResponse{
		ID:          int(item.ID),
		ProductID:   db.Int4ToInt(item.ProductID),
		ProductName: item.ProductName,
		Price:       db.MustNumericToFloat64(item.ProductPrice),
		Quantity:    int(item.Quantity),
		AddedAt:     db.TimestampToString(item.AddedAt),
	}
}

// Create a new cart for the user
func (server *Server) CreateCart(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.UserID

	cart, err := server.store.CreateCart(ctx, db.IntToInt4(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"cart_id": int(cart.ID)})
}

// Add a product to the cart
func (server *Server) AddToCart(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.UserID // Extract user ID from JWT

	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Get the authenticated user's cart
	cart, err := server.store.GetCartByUser(ctx, db.IntToInt4(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find cart"})
		return
	}

	err = server.store.AddToCart(ctx, db.AddToCartParams{
		CartID:    db.Int32ToInt4(cart.ID),
		ProductID: db.IntToInt4(req.ProductID),
		Quantity:  int32(req.Quantity),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// Retrieve the cart details
func (server *Server) GetCart(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.UserID // Extract user ID from JWT

	cart, err := server.store.GetCartByUser(ctx, db.IntToInt4(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	cartItems, err := server.store.GetCart(ctx, db.Int32ToInt4(cart.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
		return
	}

	var responses []cartItemResponse
	for _, item := range cartItems {
		responses = append(responses, newCartItemResponse(item))
	}

	ctx.JSON(http.StatusOK, gin.H{"cart_items": responses})
}

// Clear all items from the cart
func (server *Server) ClearCart(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.UserID // Extract user ID from JWT

	// Get the authenticated user's cart
	cart, err := server.store.GetCartByUser(ctx, db.IntToInt4(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find cart"})
		return
	}

	err = server.store.ClearCart(ctx, db.Int32ToInt4(cart.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
