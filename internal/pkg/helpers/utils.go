package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func MockJsonBinding(c *gin.Context, content interface{}, method string) {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
