package book

import (
	"context"
	"log"
	"net/http"

	"github.com/erizkiatama/gotu-assignment/internal/model/book"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -package=book -destination=handler_mock_test.go
type bookService interface {
	List(ctx context.Context) (book.BookResponses, error)
}

type Handler struct {
	bookSvc bookService
}

func New(bookSvc bookService) *Handler {
	return &Handler{
		bookSvc: bookSvc,
	}
}

func (h *Handler) List(c *gin.Context) {
	res, err := h.bookSvc.List(c.Request.Context())
	if err != nil {
		log.Printf("[BookHandler.List] %v", err)
		helpers.GenerateErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Response{Result: res})
}
