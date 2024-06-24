package order

import (
	"context"
	"errors"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/erizkiatama/gotu-assignment/internal/model/order"
	. "github.com/smartystreets/goconvey/convey"
)

func newMock(mockOrderRepo *MockorderReposistory) *service {
	return New(mockOrderRepo)
}

func Test_service_CreateOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	orderRepo := NewMockorderReposistory(mockCtrl)
	defer mockCtrl.Finish()

	svc := newMock(orderRepo)

	type args struct {
		userID int64
		req    order.CreateOrderRequest
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    *order.OrderResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID: 1,
				req: order.CreateOrderRequest{
					Details: []order.CreateOrderDetailRequest{
						{
							BookID: 1,
							Qty:    1,
							Price:  10000,
						},
					},
				},
			},
			mock: func(arg args) {
				orderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				orderRepo.EXPECT().BulkCreateOrderDetail(gomock.Any(), gomock.Any()).
					Return([]order.OrderDetailModel{
						{
							ID:      1,
							OrderID: 1,
							BookID:  1,
							Qty:     1,
							Price:   10000,
						},
					}, nil)
			},
			want: &order.OrderResponse{
				ID:         1,
				UserID:     1,
				TotalQty:   1,
				TotalPrice: 10000,
				Details: []order.OrderDetailResponse{
					{
						ID:     1,
						BookID: 1,
						Qty:    1,
						Price:  10000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failed to create order",
			args: args{
				userID: 1,
				req: order.CreateOrderRequest{
					Details: []order.CreateOrderDetailRequest{
						{
							BookID: 1,
							Qty:    1,
							Price:  10000,
						},
					},
				},
			},
			mock: func(arg args) {
				orderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to create order detail",
			args: args{
				userID: 1,
				req: order.CreateOrderRequest{
					Details: []order.CreateOrderDetailRequest{
						{
							BookID: 1,
							Qty:    1,
							Price:  10000,
						},
					},
				},
			},
			mock: func(arg args) {
				orderRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				orderRepo.EXPECT().BulkCreateOrderDetail(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	Convey("Test Order Service - CreateOrder", t, func() {
		for _, tt := range tests {
			tt := tt
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := svc.CreateOrder(context.Background(), tt.args.userID, tt.args.req)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(got, ShouldResemble, tt.want)
				}
			})
		}
	})
}

func Test_service_ListOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	orderRepo := NewMockorderReposistory(mockCtrl)
	defer mockCtrl.Finish()

	svc := newMock(orderRepo)

	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    []order.OrderResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID: 1,
			},
			mock: func(arg args) {
				orderRepo.EXPECT().GetAllOrder(gomock.Any(), gomock.Any()).Return([]order.OrderModel{
					{
						ID:         1,
						UserID:     1,
						TotalQty:   1,
						TotalPrice: 10000,
					},
				}, nil)
			},
			want: []order.OrderResponse{
				{
					ID:         1,
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 10000,
					Details:    nil,
				},
			},
			wantErr: false,
		},
		{
			name: "failed to list order",
			args: args{
				userID: 1,
			},
			mock: func(arg args) {
				orderRepo.EXPECT().GetAllOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	Convey("Test Order Service - ListOrder", t, func() {
		for _, tt := range tests {
			tt := tt
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := svc.ListOrder(context.Background(), tt.args.userID)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(got, ShouldResemble, tt.want)
				}
			})
		}
	})
}

func Test_service_DetailOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	orderRepo := NewMockorderReposistory(mockCtrl)
	defer mockCtrl.Finish()

	svc := newMock(orderRepo)

	type args struct {
		userID  int64
		orderID int64
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    *order.OrderResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID:  1,
				orderID: 1,
			},
			mock: func(arg args) {
				orderRepo.EXPECT().GetOrderDetail(gomock.Any(), gomock.Any()).Return([]order.OrderDetailModel{
					{
						ID:      1,
						OrderID: 1,
						BookID:  1,
						Qty:     1,
						Price:   10000,
					},
				}, nil)
			},
			want: &order.OrderResponse{
				ID:         1,
				UserID:     1,
				TotalQty:   1,
				TotalPrice: 10000,
				Details: []order.OrderDetailResponse{
					{
						ID:     1,
						BookID: 1,
						Qty:    1,
						Price:  10000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failed to get order detail",
			args: args{
				userID:  1,
				orderID: 1,
			},
			mock: func(arg args) {
				orderRepo.EXPECT().GetOrderDetail(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	Convey("Test Order Service - DetailOrder", t, func() {
		for _, tt := range tests {
			tt := tt
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := svc.DetailOrder(context.Background(), tt.args.userID, tt.args.orderID)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(got, ShouldResemble, tt.want)
				}
			})
		}
	})
}
