package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureStdout(cb func()) bytes.Buffer {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cb()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf
}

func TestNewLogger(t *testing.T) {
	assert := assert.New(t)

	buf := captureStdout(func() {
		logger := NewLogger("test-logger")

		req, _ := http.NewRequest("GET", "/test-url", nil)
		logger.Info("my message", slog.Any("request", req))
	})

	var obj struct {
		Level       string `json:"level"`
		Msg         string `json:"msg"`
		ServiceName string `json:"service_name"`
		Request     struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"request"`
	}
	json.Unmarshal(buf.Bytes(), &obj)
	assert.Equal("INFO", obj.Level)
	assert.Equal("my message", obj.Msg)
	assert.Equal("test-logger", obj.ServiceName)
	assert.Equal("GET", obj.Request.Method)
	assert.Equal("/test-url", obj.Request.Path)
}

func TestLoggerContextAttachment(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ctx := context.Background()

	ctx2 := ContextWithLogger(ctx, logger)

	logger2 := LoggerFromContext(ctx2)

	assert.Equal(logger, logger2)
}

func TestMiddlewareAddsLogContext(t *testing.T) {
	assert := assert.New(t)

	buf := captureStdout(func() {
		logger := NewLogger("test-middleware")

		handler := Middleware(logger)

		x := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := LoggerFromContext(r.Context())
			l.Error("unhandled error")
		}))

		req, _ := http.NewRequest("POST", "/path", nil)
		x.ServeHTTP(nil, req)
	})

	var obj struct {
		Level       string `json:"level"`
		Msg         string `json:"msg"`
		ServiceName string `json:"service_name"`
		TraceID     string `json:"trace_id"`
		Request     struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"request"`
	}
	log.Print(buf.String())
	json.Unmarshal(buf.Bytes(), &obj)
	assert.Equal("ERROR", obj.Level)
	assert.Equal("unhandled error", obj.Msg)
	assert.Equal("POST", obj.Request.Method)
	assert.Equal("/path", obj.Request.Path)
	assert.NotNil(obj.TraceID)
}
