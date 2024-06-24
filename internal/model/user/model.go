package user

import (
	"database/sql"
	"time"
)

type UserModel struct {
	ID        int64        `db:"id"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Name      string       `db:"name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	IsDeleted bool         `db:"is_deleted"`
}

// Requests
type (
	RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

// Responses
type (
	TokenPairResponse struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
)
