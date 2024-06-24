package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/erizkiatama/gotu-assignment/internal/model/user"
	"github.com/jmoiron/sqlx"

	. "github.com/smartystreets/goconvey/convey"
)

func newMock() (*repository, sqlmock.Sqlmock, *sql.DB) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	return New(sqlx.NewDb(db, "sqlmock")), mock, db
}

func Test_repository_Create(t *testing.T) {
	repo, mock, db := newMock()
	defer db.Close()

	type args struct {
		req user.UserModel
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    *user.UserModel
		wantErr bool
	}{
		{
			name: "error when preparing query",
			args: args{req: user.UserModel{
				Email:    "test@testing.com",
				Password: "password",
				Name:     "test",
			}},
			mock: func(args args) {
				mock.ExpectPrepare(queryCreate).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when executing query",
			args: args{req: user.UserModel{
				Email:    "test@testing.com",
				Password: "password",
				Name:     "test",
			}},
			mock: func(args args) {
				mock.ExpectPrepare(queryCreate).ExpectQuery().WithArgs(args.req.Email, args.req.Password, args.req.Name).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when scanning row",
			args: args{req: user.UserModel{
				Email:    "test@testing.com",
				Password: "password",
				Name:     "test",
			}},
			mock: func(args args) {
				mock.ExpectPrepare(queryCreate).ExpectQuery().WithArgs(args.req.Email, args.req.Password, args.req.Name).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("NaN"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{req: user.UserModel{
				Email:    "test@testing.com",
				Password: "password",
				Name:     "test",
			}},
			mock: func(args args) {
				mock.ExpectPrepare(queryCreate).ExpectQuery().WithArgs(args.req.Email, args.req.Password, args.req.Name).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			want: &user.UserModel{
				ID:       1,
				Email:    "test@testing.com",
				Password: "password",
				Name:     "test",
			},
		},
	}

	Convey("Test Create User", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := repo.Create(context.Background(), tt.args.req)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				}

				So(got, ShouldEqual, tt.want)
			})
		}
	})
}

func Test_repository_GetByEmail(t *testing.T) {
	repo, mock, db := newMock()
	defer db.Close()

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    *user.UserModel
		wantErr bool
	}{
		{
			name: "error when preparing query",
			args: args{email: "test@testing.com"},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetByEmail).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when executing query",
			args: args{email: "test@testing.com"},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetByEmail).ExpectQuery().WithArgs(args.email).WillReturnError(errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{email: "test@testing.com"},
			mock: func(args args) {
				mock.ExpectPrepare(queryGetByEmail).ExpectQuery().WithArgs(args.email).WillReturnRows(
					sqlmock.NewRows([]string{"id", "email", "password", "name"}).AddRow(1, "test@testing.com", "password", "test"),
				)
			},
			want: &user.UserModel{
				ID:       1,
				Email:    "test@testing.com",
				Password: "password",
				Name:     "test",
			},
		},
	}

	Convey("Test Get User By Email", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := repo.GetByEmail(context.Background(), tt.args.email)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				}

				So(got, ShouldResemble, tt.want)
			})
		}
	})
}
