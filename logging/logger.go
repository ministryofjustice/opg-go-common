package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Logger struct {
	serviceName string

	mu  sync.Mutex
	out *json.Encoder
}

func New(out io.Writer, serviceName string) *Logger {
	return &Logger{out: json.NewEncoder(out), serviceName: serviceName}
}

type logEvent struct {
	ServiceName   string      `json:"service_name"`
	Timestamp     time.Time   `json:"timestamp"`
	RequestMethod string      `json:"request_method,omitempty"`
	RequestURI    string      `json:"request_uri,omitempty"`
	Message       string      `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

type simpleLogEvent struct {
	ServiceName   string    `json:"service_name"`
	Timestamp     time.Time `json:"timestamp"`
	RequestMethod string    `json:"request_method,omitempty"`
	RequestURI    string    `json:"request_uri,omitempty"`
}

func (l *Logger) Print(v ...interface{}) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	_ = l.out.Encode(logEvent{
		ServiceName: l.serviceName,
		Message:     fmt.Sprint(v...),
		Timestamp:   now,
	})
}

func (l *Logger) SimplePrint(r *http.Request, err error) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	_ = l.out.Encode(simpleLogEvent{
		ServiceName:   l.serviceName,
		RequestMethod: r.Method,
		RequestURI:    r.RequestURI,
		Timestamp:     now,
	})
}

func (l *Logger) Fatal(err error) {
	l.Print(err)
	os.Exit(1)
}

// An ExpandedError will set message and data of the line logged by Request to
// Title and Data respectively.
type ExpandedError interface {
	Title() string
	Data() interface{}
}

func (l *Logger) Request(r *http.Request, err error) {
	now := time.Now()

	event := logEvent{
		ServiceName:   l.serviceName,
		RequestMethod: r.Method,
		RequestURI:    r.URL.String(),
		Timestamp:     now,
	}

	if ee, ok := err.(ExpandedError); ok {
		event.Message = ee.Title()
		event.Data = ee.Data()
	} else if err != nil {
		event.Message = err.Error()
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.out.Encode(event)
}

func (l *Logger) Response(req *http.Request, resp *http.Response, err error) {

	event := logEvent{
		ServiceName:   l.serviceName,
		RequestMethod: req.Method,
		RequestURI:    req.URL.String(),
		Timestamp:     time.Now(),
		Message:       strconv.Itoa(resp.StatusCode),
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_ = l.out.Encode(event)
}
