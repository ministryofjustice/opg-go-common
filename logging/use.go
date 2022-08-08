package logging

import "net/http"

func Use(l *Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Request(r, nil)
			next.ServeHTTP(w, r)
		})
	}
}
