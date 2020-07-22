package mux

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyHTTPResponse(t *testing.T) {
	status := http.StatusBadGateway
	headers := http.Header{"Foo": []string{"bar", "baz"}}
	body := "bang"

	resp := &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Header:     headers,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}

	rec := httptest.NewRecorder()
	copyHTTPResponse(resp, rec)
	result := rec.Result()
	assert.Equal(t, status, result.StatusCode)
	assert.Equal(t, headers, rec.Header())
	assert.Equal(t, body, rec.Body.String())
}
