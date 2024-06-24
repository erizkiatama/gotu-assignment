package book

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erizkiatama/gotu-assignment/internal/constant"
	"github.com/erizkiatama/gotu-assignment/internal/model/book"
	"github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

func newMock(mockBookSvc *MockbookService) *Handler {
	return New(mockBookSvc)
}

func Test_handler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	bookSvc := NewMockbookService(mockCtrl)
	defer mockCtrl.Finish()

	h := newMock(bookSvc)

	tests := []struct {
		name       string
		mock       func(c *gin.Context)
		wantStatus int
		want       book.BookResponses
		wantErr    bool
		err        string
	}{
		{
			name: "success",
			mock: func(c *gin.Context) {
				bookSvc.EXPECT().List(gomock.Any()).Return(book.BookResponses{
					{
						ID:          1,
						Title:       "Book 1",
						Author:      "Author 1",
						Description: "Description 1",
						Price:       150000,
					},
					{
						ID:          2,
						Title:       "Book 2",
						Author:      "Author 2",
						Description: "Description 2",
						Price:       250000,
					},
				}, nil)
			},
			wantStatus: http.StatusOK,
			want: book.BookResponses{
				{
					ID:          1,
					Title:       "Book 1",
					Author:      "Author 1",
					Description: "Description 1",
					Price:       150000,
				},
				{
					ID:          2,
					Title:       "Book 2",
					Author:      "Author 2",
					Description: "Description 2",
					Price:       250000,
				},
			},
		},
		{
			name: "error from service",
			mock: func(c *gin.Context) {
				bookSvc.EXPECT().List(gomock.Any()).Return(nil, errors.New("error from service"))
			},
			wantStatus: http.StatusInternalServerError,
			err:        constant.ErrorInternalServer,
			wantErr:    true,
		},
	}

	Convey("Test Book Handler - List", t, func() {
		for _, tt := range tests {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			Convey(tt.name, func() {
				tt.mock(c)
				h.List(c)
				So(w.Code, ShouldEqual, tt.wantStatus)

				if tt.wantErr {
					var got map[string]string
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["error"], ShouldEqual, tt.err)
				} else {
					var got map[string]book.BookResponses
					_ = json.Unmarshal(w.Body.Bytes(), &got)
					So(got["result"], ShouldResemble, tt.want)
				}

			})
		}
	})
}
