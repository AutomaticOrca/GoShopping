package server

import (
	"errors"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int   `json:"parent_id"`
}

type categoryResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id,omitempty"`
}

func newCategoryResponse(category db.Category) categoryResponse {
	return categoryResponse{
		ID:       int(category.ID),
		Name:     category.Name,
		ParentID: db.NullInt4ToPtr(category.ParentID), // 处理 NULL 值
	}
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name:     req.Name,
		ParentID: db.Int32ToNullInt4(req.ParentID),
	}

	category, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newCategoryResponse(category)
	ctx.JSON(http.StatusOK, rsp)
}

type getCategoryRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCategory(ctx *gin.Context) {
	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategoryByID(ctx, int32(req.ID))
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newCategoryResponse(category)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listCategories(ctx *gin.Context) {
	categories, err := server.store.ListCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []categoryResponse
	for _, category := range categories {
		response = append(response, newCategoryResponse(category))
	}

	ctx.JSON(http.StatusOK, response)
}

type getSubcategoriesRequest struct {
	ParentID int `uri:"parent_id" binding:"required,min=1"`
}

// 获取子分类
func (server *Server) getSubcategories(ctx *gin.Context) {
	var req getSubcategoriesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subcategories, err := server.store.GetSubcategories(ctx, db.Int32ToNullInt4(&req.ParentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []categoryResponse
	for _, category := range subcategories {
		response = append(response, newCategoryResponse(category))
	}

	ctx.JSON(http.StatusOK, response)
}

type updateCategoryRequest struct {
	ID       int    `json:"id" binding:"required,min=1"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

// 更新分类
func (server *Server) updateCategory(ctx *gin.Context) {
	var req updateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateCategoryParams{
		ID:       int32(req.ID),
		Name:     req.Name,
		ParentID: db.Int32ToNullInt4(req.ParentID),
	}

	category, err := server.store.UpdateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newCategoryResponse(category)
	ctx.JSON(http.StatusOK, rsp)
}

type deleteCategoryRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// 删除分类
func (server *Server) deleteCategory(ctx *gin.Context) {
	var req deleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteCategory(ctx, int32(req.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
