package book

import (
	"context"
	"net/http"

	"github.com/erizkiatama/gotu-assignment/internal/constant"
	"github.com/erizkiatama/gotu-assignment/internal/model/book"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
)

//go:generate mockgen -source=service.go -package=book -destination=service_mock_test.go
type bookReposistory interface {
	GetAll(ctx context.Context) (book.BookModels, error)
}

type service struct {
	bookRepo bookReposistory
}

func New(bookRepo bookReposistory) *service {
	return &service{
		bookRepo: bookRepo,
	}
}

func (s *service) List(ctx context.Context) (book.BookResponses, error) {
	books, err := s.bookRepo.GetAll(ctx)
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorListBooksFailed,
			Err:  err,
		}
	}

	res := make(book.BookResponses, len(books))
	for i, b := range books {
		res[i] = book.BookResponse{
			ID:          b.ID,
			Title:       b.Title,
			Author:      b.Author,
			Description: b.Description.String,
			Price:       b.Price,
		}
	}

	return res, nil
}
