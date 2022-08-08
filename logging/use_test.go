package logging

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	logger := New(&buf, "service")

	handler := Use(logger)(http.NotFoundHandler())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	var v logEvent
	assert.Nil(json.NewDecoder(&buf).Decode(&v))

	assert.Equal("service", v.ServiceName)
	assert.WithinDuration(time.Now(), v.Timestamp, time.Second)
	assert.Equal("GET", v.RequestMethod)
	assert.Equal("/path", v.RequestURI)
	assert.Equal("", v.Message)
	assert.Nil(v.Data)
}
