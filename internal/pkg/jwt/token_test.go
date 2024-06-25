package jwt

import (
	"testing"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateTokenPair(t *testing.T) {
	req := TokenClaim{
		Id: 123,
	}

	Convey("GenerateTokenPair", t, func() {
		Convey("should generate token pair", func() {
			config.Get().Server.SecretKey = "secret"
			tokenPair, err := GenerateTokenPair(req)
			So(err, ShouldBeNil)
			So(tokenPair, ShouldNotBeNil)

			Convey("should generate access token", func() {
				accessToken, err := validateToken(tokenPair.Access, config.Get().Server.SecretKey)
				So(err, ShouldBeNil)

				accessClaims, ok := accessToken.Claims.(jwt.MapClaims)
				So(ok, ShouldBeTrue)

				accessUserID, ok := accessClaims["user_id"].(float64)
				So(ok, ShouldBeTrue)
				So(int64(accessUserID), ShouldEqual, req.Id)
			})

			Convey("should generate refresh token", func() {
				refreshToken, err := validateToken(tokenPair.Refresh, config.Get().Server.SecretKey)
				So(err, ShouldBeNil)

				refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
				So(ok, ShouldBeTrue)

				refreshUserID, ok := refreshClaims["user_id"].(float64)
				So(ok, ShouldBeTrue)
				So(int64(refreshUserID), ShouldEqual, req.Id)
			})
		})
	})
}
func TestAuthorizeToken(t *testing.T) {
	Convey("AuthorizeToken", t, func() {
		config.Get().Server.SecretKey = "secret"
		tokenPair, err := GenerateTokenPair(TokenClaim{Id: 123})
		So(err, ShouldBeNil)
		So(tokenPair, ShouldNotBeNil)

		Convey("should return token claim for valid access token", func() {
			tokenClaim, err := AuthorizeToken(tokenPair.Access, config.Get().Server.SecretKey, false)

			So(err, ShouldBeNil)
			So(tokenClaim, ShouldNotBeNil)
			So(tokenClaim.Id, ShouldEqual, 123)
		})

		Convey("should return token claim for valid refresh token", func() {
			tokenClaim, err := AuthorizeToken(tokenPair.Refresh, config.Get().Server.SecretKey, true)

			So(err, ShouldBeNil)
			So(tokenClaim, ShouldNotBeNil)
			So(tokenClaim.Id, ShouldEqual, 123)
		})

		Convey("should return error for invalid token", func() {
			tokenClaim, err := AuthorizeToken("invalid_token", config.Get().Server.SecretKey, false)

			So(err, ShouldNotBeNil)
			So(tokenClaim, ShouldBeNil)
		})

		Convey("should return error for access token when isRefresh is true", func() {
			tokenClaim, err := AuthorizeToken(tokenPair.Access, config.Get().Server.SecretKey, true)

			So(err, ShouldNotBeNil)
			So(tokenClaim, ShouldBeNil)
		})

		Convey("should return error for expired token", func() {
			config.Get().Jwt.AccessTokenExpiryHours = -1
			expiredToken, _ := GenerateTokenPair(TokenClaim{Id: 123})
			tokenClaim, err := AuthorizeToken(expiredToken.Access, config.Get().Server.SecretKey, false)

			So(err, ShouldNotBeNil)
			So(tokenClaim, ShouldBeNil)
		})
	})
}
