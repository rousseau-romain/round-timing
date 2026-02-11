package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/rousseau-romain/round-timing/pkg/csrftoken"
)

// WithCSRFToken injects the CSRF token into the request context so templates can access it.
func WithCSRFToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), csrftoken.ContextKey, csrf.Token(r))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
