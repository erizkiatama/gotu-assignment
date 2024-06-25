package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/erizkiatama/gotu-assignment/internal/model/order"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateOrder(ctx context.Context, req order.OrderModel) (int64, error) {
	stmt, err := r.db.PreparexContext(ctx, r.db.Rebind(queryCreate))
	if err != nil {
		return 0, fmt.Errorf("[OrderRepo.CreateOrder] failed to prepare statement: %v", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	var id int64
	if err := stmt.GetContext(ctx, &id, req.UserID, req.TotalQty, req.TotalPrice); err != nil {
		return 0, fmt.Errorf("[OrderRepo.CreateOrder] failed to execute statement: %v", err)
	}

	return id, nil
}

func (r *repository) BulkCreateOrderDetail(ctx context.Context, reqs []order.OrderDetailModel) ([]order.OrderDetailModel, error) {
	var (
		values string
		args   []interface{}
	)
	for _, req := range reqs {
		values += `(?, ?, ?, ?),`
		args = append(args, req.OrderID, req.BookID, req.Qty, req.Price)
	}

	stmt, err := r.db.PreparexContext(ctx, r.db.Rebind(fmt.Sprintf(queryCreateDetail, strings.TrimSuffix(values, ","))))
	if err != nil {
		return nil, fmt.Errorf("[OrderRepo.BulkCreateOrderDetail] failed to prepare statement: %v", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryxContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("[OrderRepo.BulkCreateOrderDetail] failed to execute statement: %v", err)
	}

	i := 0
	for rows.Next() {
		if err := rows.Scan(&reqs[i].ID); err != nil {
			return nil, fmt.Errorf("[OrderRepo.BulkCreateOrderDetail] failed to scan row: %v", err)
		}
	}

	return reqs, nil
}

func (r *repository) GetAllOrder(ctx context.Context, userID int64) ([]order.OrderModel, error) {
	var res []order.OrderModel

	stmt, err := r.db.PreparexContext(ctx, r.db.Rebind(queryGetAllOrder))
	if err != nil {
		return nil, fmt.Errorf("[OrderRepo.GetAllOrder] failed to prepare statement: %v", err)
	}

	if err := stmt.SelectContext(ctx, &res, userID); err != nil {
		return nil, fmt.Errorf("[OrderRepo.GetAllOrder] failed to execute query: %v", err)
	}

	return res, nil
}

func (r *repository) GetOrderDetail(ctx context.Context, userID, orderID int64) ([]order.OrderDetailModel, error) {
	var res []order.OrderDetailModel

	stmt, err := r.db.PreparexContext(ctx, r.db.Rebind(queryGetOrderDetail))
	if err != nil {
		return nil, fmt.Errorf("[OrderRepo.GetOrderDetail] failed to prepare statement: %v", err)
	}

	if err := stmt.SelectContext(ctx, &res, userID, orderID); err != nil {
		return nil, fmt.Errorf("[OrderRepo.GetOrderDetail] failed to execute query: %v", err)
	}

	return res, nil
}
