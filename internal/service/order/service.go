package order

import (
	"context"
	"net/http"

	"github.com/erizkiatama/gotu-assignment/internal/constant"
	"github.com/erizkiatama/gotu-assignment/internal/model/order"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
)

//go:generate mockgen -source=service.go -package=order -destination=service_mock_test.go
type orderReposistory interface {
	CreateOrder(ctx context.Context, req order.OrderModel) (int64, error)
	BulkCreateOrderDetail(ctx context.Context, req []order.OrderDetailModel) ([]order.OrderDetailModel, error)
	GetAllOrder(ctx context.Context, userID int64) ([]order.OrderModel, error)
	GetOrderDetail(ctx context.Context, userID, orderID int64) ([]order.OrderDetailModel, error)
}

type service struct {
	orderRepo orderReposistory
}

func New(orderRepo orderReposistory) *service {
	return &service{
		orderRepo: orderRepo,
	}
}

func (s *service) CreateOrder(ctx context.Context, userID int64, req order.CreateOrderRequest) (*order.OrderResponse, error) {
	var (
		totalQty, totalPrice int64
		details              []order.OrderDetailModel
		detailResp           []order.OrderDetailResponse
	)

	for _, detail := range req.Details {
		totalQty += detail.Qty
		totalPrice += detail.Price
		details = append(details, order.OrderDetailModel{
			BookID: detail.BookID,
			Qty:    detail.Qty,
			Price:  detail.Price,
		})
	}

	orderID, err := s.orderRepo.CreateOrder(ctx, order.OrderModel{
		UserID:     userID,
		TotalQty:   totalQty,
		TotalPrice: totalPrice,
	})
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorCreateOrderFailed,
			Err:  err,
		}
	}

	for i, detail := range details {
		detail.OrderID = orderID
		details[i] = detail
	}

	details, err = s.orderRepo.BulkCreateOrderDetail(ctx, details)
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorCreateOrderDetailFailed,
			Err:  err,
		}
	}

	for _, detail := range details {
		detailResp = append(detailResp, order.OrderDetailResponse{
			ID:     detail.ID,
			BookID: detail.BookID,
			Qty:    detail.Qty,
			Price:  detail.Price,
		})
	}

	return &order.OrderResponse{
		ID:         orderID,
		UserID:     userID,
		TotalQty:   totalQty,
		TotalPrice: totalPrice,
		Details:    detailResp,
	}, nil
}

func (s *service) ListOrder(ctx context.Context, userID int64) ([]order.OrderResponse, error) {
	orders, err := s.orderRepo.GetAllOrder(ctx, userID)
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorGetAllOrderFailed,
			Err:  err,
		}
	}

	res := make([]order.OrderResponse, len(orders))
	for i, o := range orders {
		res[i] = order.OrderResponse{
			ID:         o.ID,
			UserID:     o.UserID,
			TotalQty:   o.TotalQty,
			TotalPrice: o.TotalPrice,
		}
	}

	return res, nil
}

func (s *service) DetailOrder(ctx context.Context, userID, orderID int64) (*order.OrderResponse, error) {
	details, err := s.orderRepo.GetOrderDetail(ctx, userID, orderID)
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorGetOrderDetailFailed,
			Err:  err,
		}
	}

	var (
		res        = make([]order.OrderDetailResponse, len(details))
		totalPrice int64
		totalQty   int64
	)
	for i, d := range details {
		totalQty += d.Qty
		totalPrice += d.Price
		res[i] = order.OrderDetailResponse{
			ID:     d.ID,
			BookID: d.BookID,
			Qty:    d.Qty,
			Price:  d.Price,
		}
	}

	return &order.OrderResponse{
		ID:         orderID,
		UserID:     userID,
		TotalQty:   totalQty,
		TotalPrice: totalPrice,
		Details:    res,
	}, nil
}
