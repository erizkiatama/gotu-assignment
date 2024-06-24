package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erizkiatama/gotu-assignment/internal/constant"
	"github.com/erizkiatama/gotu-assignment/internal/model/order"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

func newMock(orderSvc *MockorderService) *Handler {
	return New(orderSvc)
}

func Test_handler_CreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	orderSvc := NewMockorderService(mockCtrl)
	defer mockCtrl.Finish()

	h := newMock(orderSvc)

	type args struct {
		req        order.CreateOrderRequest
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		mock    func(arg args, c *gin.Context)
		want    order.OrderResponse
		wantErr bool
		err     string
	}{
		{
			name: "invalid parameters",
			args: args{
				statusCode: http.StatusBadRequest,
				req: order.CreateOrderRequest{
					Details: []order.CreateOrderDetailRequest{
						{
							BookID: 1,
							Qty:    1,
							Price:  1000,
						},
					}},
			},
			mock: func(arg args, c *gin.Context) {
				helpers.MockJsonBinding(c, map[string]interface{}{"details": "123"}, "POST")
			},
			wantErr: true,
			err:     "invalid parameters: json: cannot unmarshal string into Go struct field CreateOrderRequest.details of type []order.CreateOrderDetailRequest",
		},
		{
			name: "error from service",
			args: args{
				statusCode: http.StatusInternalServerError,
				req: order.CreateOrderRequest{
					Details: []order.CreateOrderDetailRequest{
						{
							BookID: 1,
							Qty:    1,
							Price:  1000,
						},
					},
				},
			},
			mock: func(arg args, c *gin.Context) {
				helpers.MockJsonBinding(c, arg.req, "POST")
				orderSvc.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error from service"))
			},
			wantErr: true,
			err:     constant.ErrorInternalServer,
		},
		{
			name: "success",
			args: args{
				statusCode: http.StatusCreated,
				req: order.CreateOrderRequest{
					Details: []order.CreateOrderDetailRequest{
						{
							BookID: 1,
							Qty:    1,
							Price:  1000,
						},
					},
				},
			},
			mock: func(arg args, c *gin.Context) {
				helpers.MockJsonBinding(c, arg.req, "POST")
				orderSvc.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(&order.OrderResponse{
					ID:         1,
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 1000,
					Details: []order.OrderDetailResponse{
						{
							ID:     1,
							BookID: 1,
							Qty:    1,
							Price:  1000,
						},
					},
				}, nil)
			},
			want: order.OrderResponse{
				ID:         1,
				UserID:     1,
				TotalQty:   1,
				TotalPrice: 1000,
				Details: []order.OrderDetailResponse{
					{
						ID:     1,
						BookID: 1,
						Qty:    1,
						Price:  1000,
					},
				},
			},
		},
	}

	Convey("Test Order Handler - Create Order", t, func() {
		for _, tt := range tests {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", int64(1))
			c.Request = &http.Request{
				Header: make(http.Header),
			}

			Convey(tt.name, func() {
				tt.mock(tt.args, c)
				h.CreateOrder(c)
				So(w.Code, ShouldEqual, tt.args.statusCode)

				if tt.wantErr {
					var got map[string]string
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["error"], ShouldEqual, tt.err)
				} else {
					var got map[string]order.OrderResponse
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["result"], ShouldResemble, tt.want)
				}
			})
		}
	})
}

func Test_handler_ListOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	orderSvc := NewMockorderService(mockCtrl)
	defer mockCtrl.Finish()

	h := newMock(orderSvc)

	type args struct {
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		mock    func(arg args, c *gin.Context)
		want    []order.OrderResponse
		wantErr bool
		err     string
	}{
		{
			name: "error from service",
			args: args{
				statusCode: http.StatusInternalServerError,
			},
			mock: func(arg args, c *gin.Context) {
				orderSvc.EXPECT().ListOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New("error from service"))
			},
			wantErr: true,
			err:     constant.ErrorInternalServer,
		},
		{
			name: "success",
			args: args{
				statusCode: http.StatusOK,
			},
			mock: func(arg args, c *gin.Context) {
				orderSvc.EXPECT().ListOrder(gomock.Any(), gomock.Any()).Return([]order.OrderResponse{
					{
						ID:         1,
						UserID:     1,
						TotalQty:   1,
						TotalPrice: 1000,
						Details: []order.OrderDetailResponse{
							{
								ID:     1,
								BookID: 1,
								Qty:    1,
								Price:  1000,
							},
						},
					},
				}, nil)
			},
			want: []order.OrderResponse{
				{
					ID:         1,
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 1000,
					Details: []order.OrderDetailResponse{
						{
							ID:     1,
							BookID: 1,
							Qty:    1,
							Price:  1000,
						},
					},
				},
			},
		},
	}

	Convey("Test Order Handler - List Order", t, func() {
		for _, tt := range tests {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", int64(1))
			c.Request = &http.Request{
				Header: make(http.Header),
			}

			Convey(tt.name, func() {
				tt.mock(tt.args, c)
				h.ListOrder(c)
				So(w.Code, ShouldEqual, tt.args.statusCode)

				if tt.wantErr {
					var got map[string]string
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["error"], ShouldEqual, tt.err)
				} else {
					var got map[string][]order.OrderResponse
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["result"], ShouldResemble, tt.want)
				}
			})
		}
	})
}

func Test_handler_DetailOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	orderSvc := NewMockorderService(mockCtrl)
	defer mockCtrl.Finish()

	h := newMock(orderSvc)

	type args struct {
		orderID    string
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		mock    func(arg args, c *gin.Context)
		want    order.OrderResponse
		wantErr bool
		err     string
	}{
		{
			name: "invalid parameters",
			args: args{
				statusCode: http.StatusBadRequest,
				orderID:    "abc",
			},
			mock: func(arg args, c *gin.Context) {
				c.Request = &http.Request{
					Header: make(http.Header),
				}
			},
			wantErr: true,
			err:     "invalid parameters: order_id is required",
		},
		{
			name: "error from service",
			args: args{
				statusCode: http.StatusInternalServerError,
				orderID:    "1",
			},
			mock: func(arg args, c *gin.Context) {
				c.Request = &http.Request{
					Header: make(http.Header),
				}
				orderSvc.EXPECT().DetailOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error from service"))
			},
			wantErr: true,
			err:     constant.ErrorInternalServer,
		},
		{
			name: "success",
			args: args{
				statusCode: http.StatusOK,
				orderID:    "1",
			},
			mock: func(arg args, c *gin.Context) {
				c.Request = &http.Request{
					Header: make(http.Header),
				}
				orderSvc.EXPECT().DetailOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(&order.OrderResponse{
					ID:         1,
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 1000,
					Details: []order.OrderDetailResponse{
						{
							ID:     1,
							BookID: 1,
							Qty:    1,
							Price:  1000,
						},
					},
				}, nil)
			},
			want: order.OrderResponse{
				ID:         1,
				UserID:     1,
				TotalQty:   1,
				TotalPrice: 1000,
				Details: []order.OrderDetailResponse{
					{
						ID:     1,
						BookID: 1,
						Qty:    1,
						Price:  1000,
					},
				},
			},
		},
	}

	Convey("Test Order Handler - Detail Order", t, func() {
		for _, tt := range tests {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("user_id", int64(1))
			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = append(c.Params, gin.Param{Key: "order_id", Value: tt.args.orderID})

			Convey(tt.name, func() {
				tt.mock(tt.args, c)
				h.DetailOrder(c)
				So(w.Code, ShouldEqual, tt.args.statusCode)

				if tt.wantErr {
					var got map[string]string
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["error"], ShouldEqual, tt.err)
				} else {
					var got map[string]order.OrderResponse
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["result"], ShouldResemble, tt.want)
				}
			})
		}
	})
}
