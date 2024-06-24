package middleware

import (
	"net/http"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthorizeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var res response.Response

		const bearerSchema = "Bearer "
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			res.Error = "authorization header not given"
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		tokenString := authHeader[len(bearerSchema):]
		claim, err := jwt.AuthorizeToken(tokenString, config.Get().Server.SecretKey, false)
		if err != nil {
			res.Error = err.Error()
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		c.Set("user_id", claim.Id)
		c.Next()
	}
}
