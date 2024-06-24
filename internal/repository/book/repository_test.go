package book

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/erizkiatama/gotu-assignment/internal/model/book"
	"github.com/jmoiron/sqlx"

	. "github.com/smartystreets/goconvey/convey"
)

func newMock() (*repository, sqlmock.Sqlmock, *sql.DB) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	return New(sqlx.NewDb(db, "sqlmock")), mock, db
}

func Test_repository_GetAll(t *testing.T) {
	repo, mock, db := newMock()
	defer db.Close()

	tests := []struct {
		name       string
		mock       func()
		wantResult book.BookModels
		wantErr    bool
	}{
		{
			name: "success",
			mock: func() {
				mock.ExpectPrepare(queryGetAll).ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(1, "Book 1", "Author 1").AddRow(2, "Book 2", "Author 2"),
				)
			},
			wantResult: book.BookModels{
				{ID: 1, Title: "Book 1", Author: "Author 1"},
				{ID: 2, Title: "Book 2", Author: "Author 2"},
			},
			wantErr: false,
		},
		{
			name: "error when preparing query",
			mock: func() {
				mock.ExpectPrepare(queryGetAll).WillReturnError(errors.New("error"))
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "error when executing query",
			mock: func() {
				mock.ExpectPrepare(queryGetAll).ExpectQuery().WillReturnError(errors.New("error"))
			},
			wantResult: nil,
			wantErr:    true,
		},
	}

	Convey("Test Book Repository - Get All", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock()
				gotResult, err := repo.GetAll(context.Background())
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				}

				So(gotResult, ShouldResemble, tt.wantResult)
			})
		}
	})
}
