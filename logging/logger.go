package logging

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// A Logger writes information.
//
// Deprecated: Prefer log/slog instead (which is used as the implementation
// here, if needed for an example of how it can be used).
type Logger struct {
	l *slog.Logger
}

func New(out io.Writer, serviceName string) *Logger {
	logger := slog.New(slog.
		NewJSONHandler(out, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "level" {
					return slog.Attr{}
				}

				if a.Key == "time" {
					a.Key = "timestamp"
				}

				if a.Key == "msg" {
					a.Key = "message"
				}

				return a
			},
		}).
		WithAttrs([]slog.Attr{slog.String("service_name", serviceName)}))

	return &Logger{l: logger}
}

type logEvent struct {
	ServiceName   string      `json:"service_name"`
	Timestamp     time.Time   `json:"timestamp"`
	RequestMethod string      `json:"request_method,omitempty"`
	RequestURI    string      `json:"request_uri,omitempty"`
	Message       string      `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

func (l *Logger) Print(v ...interface{}) {
	l.l.Info(fmt.Sprint(v...))
}

func (l *Logger) Fatal(err error) {
	l.l.Info(err.Error())
	os.Exit(1)
}

// An ExpandedError will set message and data of the line logged by Request to
// Title and Data respectively.
type ExpandedError interface {
	Title() string
	Data() interface{}
}

func (l *Logger) Request(r *http.Request, err error) {
	if ee, ok := err.(ExpandedError); ok {
		l.l.Info(ee.Title(),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()),
			slog.Any("data", ee.Data()))
	} else if err != nil {
		l.l.Info(err.Error(),
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	} else {
		l.l.Info("",
			slog.String("request_method", r.Method),
			slog.String("request_uri", r.URL.String()))
	}
}
