package user

import (
	"context"
	"fmt"

	"github.com/erizkiatama/gotu-assignment/internal/model/user"
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

func (r *repository) Create(ctx context.Context, req user.UserModel) (*user.UserModel, error) {
	stmt, err := r.db.PreparexContext(ctx, r.db.Rebind(queryCreate))
	if err != nil {
		return nil, fmt.Errorf("[UserRepo.Create] failed to prepare query: %v", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.QueryContext(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, fmt.Errorf("[UserRepo.Create] failed to execute query: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		if err := rows.Scan(&req.ID); err != nil {
			return nil, fmt.Errorf("[UserRepo.Create] failed to scan row: %v", err)
		}
	}

	return &req, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*user.UserModel, error) {
	var res user.UserModel

	stmt, err := r.db.PreparexContext(ctx, r.db.Rebind(queryGetByEmail))
	if err != nil {
		return nil, fmt.Errorf("[UserRepo.GetUserByEmail] failed to prepare query: %v", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	err = stmt.GetContext(ctx, &res, email)
	if err != nil {
		return nil, fmt.Errorf("[UserRepo.GetUserByEmail] failed to execute query: %w", err)
	}

	return &res, nil
}
