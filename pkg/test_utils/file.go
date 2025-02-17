package testutils

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateTestFile(t *testing.T, filename string, content string) *multipart.FileHeader {
	t.Helper()
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)

	_, err = part.Write([]byte(content))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(int64(buf.Len()))
	require.NoError(t, err)

	fileHeader := req.MultipartForm.File["file"][0]

	return fileHeader
}
