package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/erizkiatama/gotu-assignment/internal/model/user"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaim struct {
	Id int64
}

func GenerateTokenPair(req TokenClaim) (*user.TokenPairResponse, error) {
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access":  true,
		"user_id": req.Id,
		"expires": time.Now().Add(time.Duration(config.Get().Jwt.AccessTokenExpiryHours) * time.Hour).UTC().Unix(),
	})

	at, err := access.SignedString([]byte(config.Get().Server.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("error generate access token: %v", err)
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"refresh": true,
		"user_id": req.Id,
		"expires": time.Now().Add(time.Duration(config.Get().Jwt.RefreshTokenExpiryHours) * time.Hour).UTC().Unix(),
	})

	rt, err := refresh.SignedString([]byte(config.Get().Server.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("error generate refresh token: %v", err)
	}

	return &user.TokenPairResponse{
		Access:  at,
		Refresh: rt,
	}, nil
}

func validateToken(encodedToken, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}
		return []byte(secretKey), nil
	})
}

func AuthorizeToken(tokenString, secretKey string, isRefresh bool) (*TokenClaim, error) {
	token, err := validateToken(tokenString, secretKey)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if isRefresh && !claims["refresh"].(bool) {
			return nil, errors.New("token was not a refresh token")
		}

		if int64(claims["expires"].(float64)) < time.Now().UTC().Unix() {
			return nil, errors.New("token has expired")
		}

		return &TokenClaim{
			Id: int64(claims["user_id"].(float64)),
		}, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
