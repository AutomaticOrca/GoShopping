package server

import (
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	util "github.com/AutomaticOrca/GoShopping/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createUserRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:    int(user.ID),
		Email: user.Email,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUserByEmail(ctx, req.Email)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}
