package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/erizkiatama/gotu-assignment/internal/constant"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/erizkiatama/gotu-assignment/internal/model/user"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/helpers"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -package=user -destination=service_mock_test.go
type userReposistory interface {
	Create(ctx context.Context, req user.UserModel) (*user.UserModel, error)
	GetByEmail(ctx context.Context, email string) (*user.UserModel, error)
}

type service struct {
	userRepo userReposistory
}

func New(userRepo userReposistory) *service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) Register(ctx context.Context, req user.RegisterRequest) (*user.TokenPairResponse, error) {
	hashedPassword, err := helpers.EncryptPassword([]byte(req.Password))
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorInternalServer,
			Err:  err,
		}
	}

	newUser := user.UserModel{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}
	res, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, &response.ServiceError{
				Code: http.StatusConflict,
				Msg:  constant.ErrorUserAlreadyExists,
				Err:  err,
			}
		} else {
			return nil, &response.ServiceError{
				Code: http.StatusInternalServerError,
				Msg:  constant.ErrorCreateUserFailed,
				Err:  err,
			}
		}
	}

	tokenPair, err := jwt.GenerateTokenPair(jwt.TokenClaim{
		Id: res.ID,
	})
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorGenerateToken,
			Err:  fmt.Errorf("[UserSvc.Register] failed to generate token: %v", err),
		}
	}

	return tokenPair, nil
}

func (s *service) Login(ctx context.Context, req user.LoginRequest) (*user.TokenPairResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ServiceError{
				Code: http.StatusNotFound,
				Msg:  constant.ErrorUserNotFound,
				Err:  err,
			}
		}
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorGetUserFailed,
			Err:  err,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusBadRequest,
			Msg:  constant.ErrorPasswordNotMatch,
			Err:  fmt.Errorf("[UserSvc.Login] failed to compare password: %v", err),
		}
	}

	tokenPair, err := jwt.GenerateTokenPair(jwt.TokenClaim{Id: user.ID})
	if err != nil {
		return nil, &response.ServiceError{
			Code: http.StatusInternalServerError,
			Msg:  constant.ErrorGenerateToken,
			Err:  fmt.Errorf("[UserSvc.Login] failed to generate token: %v", err),
		}
	}

	return tokenPair, nil

}
