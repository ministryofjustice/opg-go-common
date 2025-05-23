package securityheaders

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	assert := assert.New(t)

	handler := Use(http.NotFoundHandler())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal("default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'self'; upgrade-insecure-requests; block-all-mixed-content;", resp.Header.Get("Content-Security-Policy"))
	assert.Equal("same-origin", resp.Header.Get("Referrer-Policy"))
	assert.Equal("max-age=31536000; includeSubDomains; preload", resp.Header.Get("Strict-Transport-Security"))
	assert.Equal("nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal("SAMEORIGIN", resp.Header.Get("X-Frame-Options"))
	assert.Equal("1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.Equal("no-store", resp.Header.Get("Cache-control"))
	assert.Equal("no-cache", resp.Header.Get("Pragma"))
}
