package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/erizkiatama/gotu-assignment/internal/model/user"
	. "github.com/smartystreets/goconvey/convey"
)

func newMock(mockUserRepo *MockuserReposistory) *service {
	return New(mockUserRepo)
}

func Test_service_Register(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := NewMockuserReposistory(mockCtrl)
	defer mockCtrl.Finish()

	svc := newMock(userRepo)

	type args struct {
		req user.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    *user.TokenPairResponse
		wantErr bool
	}{
		{
			name: "failed to encrypt password - string more than 72 bytes",
			args: args{
				req: user.RegisterRequest{
					Email:    "test@testing.com",
					Password: "xBDVuEZUXNapHocEMCwoCvKMkhKwNXoDPwFzEmuZpzycQgubTFJcGHvgdNpqdwTuYkECqooJWDe",
					Name:     "test",
				},
			},
			mock:    func(a args) {},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to create user - duplicate key error",
			args: args{
				req: user.RegisterRequest{
					Email:    "test@testing.com",
					Password: "password",
					Name:     "test",
				},
			},
			mock: func(arg args) {
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("duplicate key error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to create user - internal server error",
			args: args{
				req: user.RegisterRequest{
					Email:    "test@testing.com",
					Password: "password",
					Name:     "test",
				},
			},
			mock: func(arg args) {
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "succcess",
			args: args{
				req: user.RegisterRequest{
					Email:    "test@testing.com",
					Password: "password",
					Name:     "test",
				},
			},
			mock: func(arg args) {
				config.Get().Server.SecretKey = "testing"
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&user.UserModel{ID: 1}, nil)
			},
		},
	}

	Convey("Test User Service - Register", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := svc.Register(context.Background(), tt.args.req)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				} else {
					So(got, ShouldNotBeNil)
				}
			})
		}
	})
}

func Test_service_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	userRepo := NewMockuserReposistory(mockCtrl)
	defer mockCtrl.Finish()

	svc := newMock(userRepo)

	type args struct {
		req user.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args)
		want    *user.TokenPairResponse
		wantErr bool
	}{
		{
			name: "user not found",
			args: args{
				req: user.LoginRequest{
					Email:    "test@testing.com",
					Password: "password",
				},
			},
			mock: func(arg args) {
				userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to get user",
			args: args{
				req: user.LoginRequest{
					Email:    "test@testing.com",
					Password: "password",
				},
			},
			mock: func(arg args) {
				userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "password not match",
			args: args{
				req: user.LoginRequest{
					Email:    "test@testing.com",
					Password: "password",
				},
			},
			mock: func(arg args) {
				userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&user.UserModel{Password: "hashed_password"}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				req: user.LoginRequest{
					Email:    "test@testing.com",
					Password: "password",
				},
			},
			mock: func(arg args) {
				userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&user.UserModel{ID: 1, Password: "$2a$10$zxOrik5iLL4LBOXNRzCTY.5QUBFFTRbKxpemQbygAN6nouR7G3CU6"}, nil)
			},
			want:    &user.TokenPairResponse{},
			wantErr: false,
		},
	}

	Convey("Test User Service - Login", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock(tt.args)
				got, err := svc.Login(context.Background(), tt.args.req)
				if tt.wantErr {
					So(err, ShouldNotBeNil)
				} else {
					So(got, ShouldNotBeNil)
				}
			})
		}
	})
}
