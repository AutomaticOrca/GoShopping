package server

import (
	"net/http"
	"strconv"

	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

type productResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
	CategoryID  int     `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Convert db.Product to productResponse
func newProductResponse(product db.Product) productResponse {
	return productResponse{
		ID:          int(product.ID),
		Name:        product.Name,
		Description: db.TextToString(product.Description),
		Price:       db.NumericToFloat64(product.Price),
		Stock:       int(product.Stock),
		ImageURL:    db.TextToString(product.ImageUrl),
		CategoryID:  int(product.CategoryID.Int32),
		CreatedAt:   db.TimestampToString(product.CreatedAt),
		UpdatedAt:   db.TimestampToString(product.UpdatedAt),
	}
}

// Create Product
func (server *Server) createProduct(ctx *gin.Context) {
	var req db.CreateProductParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := server.store.CreateProduct(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newProductResponse(product))
}

// Update Product
func (server *Server) updateProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req db.UpdateProductParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = int32(id) // Convert int64 → int32

	product, err := server.store.UpdateProduct(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newProductResponse(product))
}

// Delete Product
func (server *Server) deleteProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	err = server.store.DeleteProduct(ctx, int32(id)) // Convert int64 → int32
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// Get Product by ID
func (server *Server) getProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := server.store.GetProductByID(ctx, int32(id)) // Convert int64 → int32
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, newProductResponse(product))
}

// List Products (Paginated)
func (server *Server) listProducts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	products, err := server.store.ListProducts(ctx, db.ListProductsParams{
		Limit:  int32(pageSize), // Convert int64 → int32
		Offset: int32((page - 1) * pageSize),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []productResponse
	for _, product := range products {
		responses = append(responses, newProductResponse(product))
	}

	ctx.JSON(http.StatusOK, gin.H{"products": responses})
}
