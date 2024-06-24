package book

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/erizkiatama/gotu-assignment/internal/model/book"
	. "github.com/smartystreets/goconvey/convey"
	gomock "go.uber.org/mock/gomock"
)

func newMock(mockBookRepo *MockbookReposistory) *service {
	return New(mockBookRepo)
}

func Test_service_List(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	bookRepo := NewMockbookReposistory(mockCtrl)
	defer mockCtrl.Finish()

	svc := newMock(bookRepo)

	ctx := context.Background()

	tests := []struct {
		name       string
		mock       func()
		wantResult book.BookResponses
		wantErr    bool
	}{
		{
			name: "success",
			mock: func() {
				bookRepo.EXPECT().GetAll(ctx).Return([]book.BookModel{
					{
						ID:          1,
						Title:       "Book 1",
						Author:      "Author 1",
						Description: sql.NullString{String: "Description 1", Valid: true},
						Price:       150000,
					},
					{
						ID:          2,
						Title:       "Book 2",
						Author:      "Author 2",
						Description: sql.NullString{String: "Description 2", Valid: true},
						Price:       250000,
					},
				}, nil)
			},
			wantResult: book.BookResponses{
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
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				bookRepo.EXPECT().GetAll(ctx).Return(nil, errors.New("database error"))
			},
			wantResult: nil,
			wantErr:    true,
		},
	}

	Convey("Test Book Service - List", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				tt.mock()
				gotResult, gotErr := svc.List(ctx)
				if tt.wantErr {
					So(gotErr, ShouldNotBeNil)
				} else {
					So(gotErr, ShouldBeNil)
					So(gotResult, ShouldEqual, tt.wantResult)
				}
			})
		}
	})
}
