package helpers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
)

func TestEncryptPassword(t *testing.T) {
	Convey("Success", t, func() {
		password := []byte("password123")
		hash, err := EncryptPassword(password)

		So(err, ShouldBeNil)
		So(hash, ShouldNotBeEmpty)

		err = bcrypt.CompareHashAndPassword([]byte(hash), password)
		So(err, ShouldBeNil)
	})

	Convey("Error", t, func() {
		password := []byte("KiEdNLDaVfxMGUEniAQfcjnuEcXDPHVHkjaAJDKxKpguVaBjwvqdhVuywYXmmQBMPWEfDrCJa")
		hash, err := EncryptPassword(password)

		So(err, ShouldNotBeNil)
		So(hash, ShouldBeEmpty)
	})
}
func TestGenerateErrorResponse(t *testing.T) {
	Convey("When given a service error", t, func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		err := &response.ServiceError{
			Code: http.StatusBadRequest,
			Msg:  "Bad request",
		}

		GenerateErrorResponse(c, err)

		Convey("Should return the service error as JSON response", func() {
			resp := w.Body.String()
			expectedResp := `{"error":"Bad request"}`

			So(resp, ShouldEqual, expectedResp)
		})
	})

	Convey("When given a non-service error", t, func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		err := errors.New("Some error")

		GenerateErrorResponse(c, err)

		Convey("Should return the internal server error as JSON response", func() {
			resp := w.Body.String()
			expectedResp := `{"error":"internal server error"}`

			So(resp, ShouldEqual, expectedResp)
		})
	})
}
func TestMockJsonBinding(t *testing.T) {
	Convey("When given a valid content and method", t, func() {
		content := struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{
			Name: "John Doe",
			Age:  30,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		MockJsonBinding(c, content, http.MethodPost)

		Convey("Should set the request method and content type", func() {
			So(c.Request.Method, ShouldEqual, http.MethodPost)
			So(c.Request.Header.Get("Content-Type"), ShouldEqual, "application/json")
		})

		Convey("Should set the request body with the JSON content", func() {
			body, _ := io.ReadAll(c.Request.Body)

			expectedBody := `{"name":"John Doe","age":30}`
			So(string(body), ShouldEqual, expectedBody)
		})
	})
}
