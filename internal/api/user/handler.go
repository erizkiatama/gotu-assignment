package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/erizkiatama/gotu-assignment/internal/model/user"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -package=user -destination=handler_mock_test.go
type userService interface {
	Register(ctx context.Context, req user.RegisterRequest) (*user.TokenPairResponse, error)
	Login(ctx context.Context, req user.LoginRequest) (*user.TokenPairResponse, error)
}

type Handler struct {
	userSvc userService
}

func New(userSvc userService) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req user.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Error: fmt.Sprintf("invalid parameters: %s", err.Error()),
		})
		return
	}

	res, err := h.userSvc.Register(c.Request.Context(), req)
	if err != nil {
		log.Printf("[UserHandler.Register] %v", err)
		helpers.GenerateErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Response{Result: res})
}

func (h *Handler) Login(c *gin.Context) {
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Error: fmt.Sprintf("invalid parameters: %s", err.Error()),
		})
		return
	}

	res, err := h.userSvc.Login(c.Request.Context(), req)
	if err != nil {
		log.Printf("[UserHandler.Login] %v", err)
		helpers.GenerateErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Response{Result: res})
}
