package book

import (
	"database/sql"
	"time"
)

type (
	BookModel struct {
		ID          int64          `db:"id"`
		Title       string         `db:"title"`
		Author      string         `db:"author"`
		Description sql.NullString `db:"description"`
		Price       int64          `db:"price"`
		CreatedAt   time.Time      `db:"created_at"`
		UpdatedAt   sql.NullTime   `db:"updated_at"`
		IsDeleted   bool           `db:"is_deleted"`
	}

	BookModels []BookModel
)

// Requests

// Responses
type (
	BookResponse struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Price       int64  `json:"price"`
	}

	BookResponses []BookResponse
)
