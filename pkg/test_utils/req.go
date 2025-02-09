package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewJSONRequest(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(method, url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}
