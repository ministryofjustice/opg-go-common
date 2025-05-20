package securityheaders

import "net/http"

// Use applies common security headers to a handler.
func Use(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'self'; upgrade-insecure-requests; block-all-mixed-content;")
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
