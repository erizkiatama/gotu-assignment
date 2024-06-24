package helpers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/erizkiatama/gotu-assignment/internal/constant"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error encrypting password: %v", err)
	}

	return string(hash), nil
}

func GenerateErrorResponse(c *gin.Context, err error) {
	var svcErr *response.ServiceError

	if errors.As(err, &svcErr) {
		c.JSON(svcErr.Code, response.Response{
			Error: svcErr.Msg,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, response.Response{
		Error: constant.ErrorInternalServer,
	})
}
