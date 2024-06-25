package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/erizkiatama/gotu-assignment/internal/model/order"
	"github.com/jmoiron/sqlx"

	. "github.com/smartystreets/goconvey/convey"
)

func newMock() (*repository, sqlmock.Sqlmock, *sql.DB) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	return New(sqlx.NewDb(db, "sqlmock")), mock, db
}

func Test_repository_CreateOrder(t *testing.T) {
	repo, mock, db := newMock()
	defer db.Close()

	type args struct {
		req order.OrderModel
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    int64
		wantErr bool
	}{
		{
			name: "error when preparing query",
			args: args{
				req: order.OrderModel{
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 10000,
				}},
			mock: func(args args) {
				mock.ExpectPrepare(queryCreate).WillReturnError(errors.New("error"))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error when executing query",
			args: args{
				req: order.OrderModel{
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 10000,
				}},
			mock: func(args args) {
				mock.ExpectPrepare(queryCreate).ExpectQuery().
					WithArgs(args.req.UserID, args.req.TotalQty, args.req.TotalPrice).WillReturnError(errors.New("error"))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				req: order.OrderModel{
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 10000,
				}},
			mock: func(args args) {
				mock.ExpectPrepare(queryCreate).ExpectQuery().
					WithArgs(args.req.UserID, args.req.TotalQty, args.req.TotalPrice).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			want:    1,
			wantErr: false,
		},
	}

	Convey("Test Order Repository - Create Order", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := repo.CreateOrder(context.Background(), tt.args.req)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				}
				So(got, ShouldEqual, tt.want)
			})
		}
	})
}

func Test_repository_BulkCreateOrderDetail(t *testing.T) {
	repo, mock, db := newMock()
	defer db.Close()

	type args struct {
		reqs []order.OrderDetailModel
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    []order.OrderDetailModel
		wantErr bool
	}{
		{
			name: "error when preparing query",
			args: args{
				reqs: []order.OrderDetailModel{
					{
						OrderID: 1,
						BookID:  1,
						Qty:     1,
						Price:   10000,
					},
				},
			},
			mock: func(args args) {
				mock.ExpectPrepare(fmt.Sprintf(queryCreateDetail, "(?, ?, ?, ?)")).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when executing query",
			args: args{
				reqs: []order.OrderDetailModel{
					{
						OrderID: 1,
						BookID:  1,
						Qty:     1,
						Price:   10000,
					},
				},
			},
			mock: func(args args) {
				mock.ExpectPrepare(fmt.Sprintf(queryCreateDetail, "(?, ?, ?, ?)")).ExpectQuery().
					WithArgs(args.reqs[0].OrderID, args.reqs[0].BookID, args.reqs[0].Qty, args.reqs[0].Price).
					WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when scanning row",
			args: args{
				reqs: []order.OrderDetailModel{
					{
						OrderID: 1,
						BookID:  1,
						Qty:     1,
						Price:   10000,
					},
				},
			},
			mock: func(args args) {
				mock.ExpectPrepare(fmt.Sprintf(queryCreateDetail, "(?, ?, ?, ?)")).ExpectQuery().
					WithArgs(args.reqs[0].OrderID, args.reqs[0].BookID, args.reqs[0].Qty, args.reqs[0].Price).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("NaN"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				reqs: []order.OrderDetailModel{
					{
						OrderID: 1,
						BookID:  1,
						Qty:     1,
						Price:   10000,
					},
				},
			},
			mock: func(args args) {
				mock.ExpectPrepare(fmt.Sprintf(queryCreateDetail, "(?, ?, ?, ?)")).ExpectQuery().
					WithArgs(args.reqs[0].OrderID, args.reqs[0].BookID, args.reqs[0].Qty, args.reqs[0].Price).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			want: []order.OrderDetailModel{
				{
					ID:      1,
					OrderID: 1,
					BookID:  1,
					Qty:     1,
					Price:   10000,
				},
			},
			wantErr: false,
		},
	}

	Convey("Test Order Repository - Bulk Create Order Detail", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := repo.BulkCreateOrderDetail(context.Background(), tt.args.reqs)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				}
				So(got, ShouldResemble, tt.want)
			})
		}
	})
}

func Test_repository_GetAllOrder(t *testing.T) {
	repo, mock, db := newMock()
	defer db.Close()

	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    []order.OrderModel
		wantErr bool
	}{
		{
			name: "error when preparing query",
			args: args{userID: 1},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetAllOrder).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when executing query",
			args: args{userID: 1},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetAllOrder).ExpectQuery().WithArgs(args.userID).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{userID: 1},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetAllOrder).ExpectQuery().WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "total_quantity", "total_price"}).
						AddRow(1, 1, 1, 10000))
			},
			want: []order.OrderModel{
				{
					ID:         1,
					UserID:     1,
					TotalQty:   1,
					TotalPrice: 10000,
				},
			},
			wantErr: false,
		},
	}

	Convey("Test Order Repository - Get All Order", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := repo.GetAllOrder(context.Background(), tt.args.userID)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				}
				So(got, ShouldResemble, tt.want)
			})
		}
	})
}

func Test_repository_GetOrderDetail(t *testing.T) {
	repo, mock, db := newMock()
	defer db.Close()

	type args struct {
		orderID int64
		userID  int64
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    []order.OrderDetailModel
		wantErr bool
	}{
		{
			name: "error when preparing query",
			args: args{orderID: 1},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetOrderDetail).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when executing query",
			args: args{orderID: 1},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetOrderDetail).ExpectQuery().WithArgs(args.userID, args.orderID).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{orderID: 1},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetOrderDetail).ExpectQuery().WithArgs(args.userID, args.orderID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "book_id", "quantity", "price"}).
						AddRow(1, 1, 1, 1, 10000))
			},
			want: []order.OrderDetailModel{
				{
					ID:      1,
					OrderID: 1,
					BookID:  1,
					Qty:     1,
					Price:   10000,
				},
			},
			wantErr: false,
		},
	}

	Convey("Test Order Repository - Get Order Detail", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := repo.GetOrderDetail(context.Background(), tt.args.userID, tt.args.orderID)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				}
				So(got, ShouldResemble, tt.want)
			})
		}
	})
}
