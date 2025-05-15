package securityheaders

import "net/http"

// Use applies common security headers to a handler.
func Use(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; object-src 'none'; connect-src 'self'; font-src 'self'; frame-ancestors 'none';")
		w.Header().Add("Referrer-Policy", "same-origin")
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Frame-Options", "SAMEORIGIN")
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("Cache-control", "no-store")
		w.Header().Add("Pragma", "no-cache")

		h.ServeHTTP(w, r)
	}
}
