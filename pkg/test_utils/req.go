package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewJSONRequest(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(method, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func NewMultipartRequest(t *testing.T, method, url string, data map[string]any) *http.Request {
	t.Helper()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer writer.Close()

	for fieldname, value := range data {
		switch value := value.(type) {
		case string:
			require.NoError(t, writer.WriteField(fieldname, value))
		case []byte:
			part, err := writer.CreateFormFile(fieldname, fmt.Sprintf("%s.jpg", fieldname))
			require.NoError(t, err)
			_, err = part.Write(value)
			require.NoError(t, err)
		}
	}

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
