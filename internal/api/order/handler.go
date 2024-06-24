package order

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erizkiatama/gotu-assignment/internal/model/order"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -package=order -destination=handler_mock_test.go
type orderService interface {
	CreateOrder(ctx context.Context, userID int64, req order.CreateOrderRequest) (*order.OrderResponse, error)
	ListOrder(ctx context.Context, userID int64) ([]order.OrderResponse, error)
	DetailOrder(ctx context.Context, userID, orderID int64) (*order.OrderResponse, error)
}

type Handler struct {
	orderSvc orderService
}

func New(orderSvc orderService) *Handler {
	return &Handler{
		orderSvc: orderSvc,
	}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req order.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Error: fmt.Sprintf("invalid parameters: %s", err.Error()),
		})
		return
	}

	userID, _ := c.Get("user_id")
	res, err := h.orderSvc.CreateOrder(c.Request.Context(), userID.(int64), req)
	if err != nil {
		log.Printf("[OrderHandler.CreateOrder] %v", err)
		helpers.GenerateErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.Response{Result: res})
}

func (h *Handler) ListOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	res, err := h.orderSvc.ListOrder(c.Request.Context(), userID.(int64))
	if err != nil {
		log.Printf("[OrderHandler.ListOrder] %v", err)
		helpers.GenerateErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Response{Result: res})
}

func (h *Handler) DetailOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderIDStr := c.Param("order_id")

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if orderID == 0 || err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Error: "invalid parameters: order_id is required",
		})
		return
	}

	res, err := h.orderSvc.DetailOrder(c.Request.Context(), userID.(int64), orderID)
	if err != nil {
		log.Printf("[OrderHandler.DetailOrder] %v", err)
		helpers.GenerateErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Response{Result: res})
}
