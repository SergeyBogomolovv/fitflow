package content_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	contentHandler "github.com/SergeyBogomolovv/fitflow/internal/delivery/http/content"
	"github.com/SergeyBogomolovv/fitflow/internal/delivery/http/content/mocks"
	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContentHandler_GenerateContent(t *testing.T) {
	type MockBehavior func(repo *mocks.ContentService, theme string)

	testCases := []struct {
		name           string
		theme          string
		mockBehavior   MockBehavior
		wantStatusCode int
		wantBody       string
	}{
		{
			name:  "success",
			theme: "test_theme",
			mockBehavior: func(svc *mocks.ContentService, theme string) {
				svc.EXPECT().GenerateContent(mock.Anything, theme).Return("test_content", nil).Once()
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `{"status":"success","content":"test_content"}` + "\n",
		},
		{
			name:           "no theme",
			theme:          "",
			mockBehavior:   func(svc *mocks.ContentService, theme string) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"status":"error","code":400,"message":"theme is required"}` + "\n",
		},
		{
			name:  "error",
			theme: "test_theme",
			mockBehavior: func(svc *mocks.ContentService, theme string) {
				svc.EXPECT().GenerateContent(mock.Anything, theme).Return("", fmt.Errorf("error")).Once()
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       `{"status":"error","code":500,"message":"failed to generate content"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contentSvc := mocks.NewContentService(t)
			tc.mockBehavior(contentSvc, tc.theme)

			handler := contentHandler.New(testutils.NewTestLogger(), contentSvc)

			rec := httptest.NewRecorder()
			url := fmt.Sprintf("/content/generate?theme=%s", tc.theme)
			req := testutils.NewJSONRequest(t, http.MethodGet, url, nil)
			handler.HandleGenerateContent(rec, req)

			assert.Equal(t, tc.wantStatusCode, rec.Code)
			assert.Equal(t, tc.wantBody, rec.Body.String())
		})
	}
}

func TestContentHandler_CreatePost(t *testing.T) {
	type args struct {
		content   string
		audience  domain.UserLvl
		withImage bool
	}

	type MockBehavior func(svc *mocks.ContentService, args args)

	testCases := []struct {
		name           string
		args           args
		mockBehavior   MockBehavior
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "success",
			args: args{content: "test content", audience: domain.UserLvlDefault, withImage: true},
			mockBehavior: func(svc *mocks.ContentService, args args) {
				svc.EXPECT().CreatePost(mock.Anything, mock.Anything).
					Return(domain.Post{
						ID:       1,
						Content:  args.content,
						Audience: args.audience,
						Images:   []string{"http://image.ru"},
					}, nil).Once()
			},
			wantStatusCode: http.StatusCreated,
			wantBody:       `{"id":1,"content":"test content","audience":"default","images":["http://image.ru"]}` + "\n",
		},
		{
			name:           "without image",
			args:           args{content: "test content", audience: domain.UserLvlDefault, withImage: false},
			mockBehavior:   func(svc *mocks.ContentService, args args) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"status":"error","code":400,"message":"invalid payload"}` + "\n",
		},
		{
			name:           "invalid audience",
			args:           args{content: "test content", audience: domain.UserLvl("sfsf"), withImage: true},
			mockBehavior:   func(svc *mocks.ContentService, args args) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"status":"error","code":400,"message":"invalid payload"}` + "\n",
		},
		{
			name:           "no content",
			args:           args{content: "", audience: domain.UserLvlDefault, withImage: true},
			mockBehavior:   func(svc *mocks.ContentService, args args) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"status":"error","code":400,"message":"invalid payload"}` + "\n",
		},
		{
			name: "error",
			args: args{content: "test content", audience: domain.UserLvlDefault, withImage: true},
			mockBehavior: func(svc *mocks.ContentService, args args) {
				svc.EXPECT().CreatePost(mock.Anything, mock.Anything).Return(domain.Post{}, assert.AnError).Once()
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       `{"status":"error","code":500,"message":"failed to create post"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contentSvc := mocks.NewContentService(t)
			tc.mockBehavior(contentSvc, tc.args)

			handler := contentHandler.New(testutils.NewTestLogger(), contentSvc)
			rec := httptest.NewRecorder()
			body := map[string]any{
				"content":  tc.args.content,
				"audience": string(tc.args.audience),
			}
			if tc.args.withImage {
				body["images"] = []byte("image_data")
			}

			req := testutils.NewMultipartRequest(t, http.MethodPost, "/content/post", body)
			handler.HandleCreatePost(rec, req)

			assert.Equal(t, tc.wantStatusCode, rec.Code)
			assert.Equal(t, tc.wantBody, rec.Body.String())
		})
	}
}
