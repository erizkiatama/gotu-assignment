package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erizkiatama/gotu-assignment/internal/constant"
	"github.com/erizkiatama/gotu-assignment/internal/model/response"
	"github.com/erizkiatama/gotu-assignment/internal/model/user"
	"github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

func newMock(userSvc *MockuserService) *Handler {
	return New(userSvc)
}

func mockJsonBinding(c *gin.Context, content interface{}, method string) {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func Test_handler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	userSvc := NewMockuserService(mockCtrl)
	defer mockCtrl.Finish()

	h := newMock(userSvc)

	type args struct {
		req        user.RegisterRequest
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		mock    func(arg args, c *gin.Context)
		want    response.Response
		wantErr bool
	}{
		{
			name: "invalid parameters",
			args: args{
				statusCode: http.StatusBadRequest,
				req: user.RegisterRequest{
					Email:    "test@testing.com",
					Password: "password",
					Name:     "test",
				}},
			mock: func(arg args, c *gin.Context) {
				mockJsonBinding(c, map[string]interface{}{"email": 123}, "POST")
			},
			want: response.Response{
				Error: "invalid parameters: json: cannot unmarshal number into Go struct field RegisterRequest.email of type string",
			},
		},
		{
			name: "error from service",
			args: args{
				statusCode: http.StatusConflict,
				req: user.RegisterRequest{
					Email:    "test@testing.com",
					Password: "password",
					Name:     "test",
				}},
			mock: func(arg args, c *gin.Context) {
				mockJsonBinding(c, arg.req, "POST")
				userSvc.EXPECT().Register(gomock.Any(), arg.req).Return(nil, &response.ServiceError{
					Code: http.StatusConflict,
					Msg:  constant.ErrorUserAlreadyExists,
					Err:  errors.New("error user already exists"),
				})
			},
			want: response.Response{
				Error: constant.ErrorUserAlreadyExists,
			},
		},
		{
			name: "success",
			args: args{
				statusCode: http.StatusCreated,
				req: user.RegisterRequest{
					Email:    "test@testing.com",
					Password: "password",
					Name:     "test",
				}},
			mock: func(arg args, c *gin.Context) {
				mockJsonBinding(c, arg.req, "POST")
				userSvc.EXPECT().Register(gomock.Any(), arg.req).Return(&user.TokenPairResponse{
					Access:  "access_token",
					Refresh: "refresh_token",
				}, nil)
			},
			want: response.Response{
				Result: map[string]interface{}{
					"access":  "access_token",
					"refresh": "refresh_token",
				},
			},
		},
	}
	Convey("Test User Handler - Register", t, func() {
		for _, tt := range tests {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			Convey(tt.name, func() {
				tt.mock(tt.args, c)
				h.Register(c)
				So(w.Code, ShouldEqual, tt.args.statusCode)

				var got response.Response
				_ = json.Unmarshal(w.Body.Bytes(), &got)
				So(got, ShouldEqual, tt.want)
			})
		}
	})
}
func Test_handler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	userSvc := NewMockuserService(mockCtrl)
	defer mockCtrl.Finish()

	h := newMock(userSvc)

	type args struct {
		req        user.LoginRequest
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		mock    func(arg args, c *gin.Context)
		want    response.Response
		wantErr bool
	}{
		{
			name: "invalid parameters",
			args: args{
				statusCode: http.StatusBadRequest,
				req: user.LoginRequest{
					Email:    "test@testing.com",
					Password: "password",
				},
			},
			mock: func(arg args, c *gin.Context) {
				mockJsonBinding(c, map[string]interface{}{"email": 123}, "POST")
			},
			want: response.Response{
				Error: "invalid parameters: json: cannot unmarshal number into Go struct field LoginRequest.email of type string",
			},
		},
		{
			name: "error from service",
			args: args{
				statusCode: http.StatusInternalServerError,
				req: user.LoginRequest{
					Email:    "test@testing.com",
					Password: "password",
				},
			},
			mock: func(arg args, c *gin.Context) {
				mockJsonBinding(c, arg.req, "POST")
				userSvc.EXPECT().Login(gomock.Any(), arg.req).Return(nil, errors.New("error from service"))
			},
			want: response.Response{
				Error: constant.ErrorInternalServer,
			},
		},
		{
			name: "success",
			args: args{
				statusCode: http.StatusOK,
				req: user.LoginRequest{
					Email:    "test@testing.com",
					Password: "password",
				},
			},
			mock: func(arg args, c *gin.Context) {
				mockJsonBinding(c, arg.req, "POST")
				userSvc.EXPECT().Login(gomock.Any(), arg.req).Return(&user.TokenPairResponse{
					Access:  "access_token",
					Refresh: "refresh_token",
				}, nil)
			},
			want: response.Response{
				Result: map[string]interface{}{
					"access":  "access_token",
					"refresh": "refresh_token",
				},
			},
		},
	}

	Convey("Test User Handler - Login", t, func() {
		for _, tt := range tests {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			Convey(tt.name, func() {
				tt.mock(tt.args, c)
				h.Login(c)
				So(w.Code, ShouldEqual, tt.args.statusCode)

				var got response.Response
				_ = json.Unmarshal(w.Body.Bytes(), &got)
				So(got, ShouldEqual, tt.want)
			})
		}
	})
}
