package book

import (
	"context"
	"fmt"

	"github.com/erizkiatama/gotu-assignment/internal/model/book"
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

func (r *repository) GetAll(ctx context.Context) (book.BookModels, error) {
	var res book.BookModels

	stmt, err := r.db.PreparexContext(ctx, r.db.Rebind(queryGetAll))
	if err != nil {
		return nil, fmt.Errorf("[BookRepo.GetAll] failed to prepare query: %v", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	err = stmt.SelectContext(ctx, &res)
	if err != nil {
		return nil, fmt.Errorf("[BookRepo.GetAll] failed to execute query: %v", err)
	}

	return res, nil
}
