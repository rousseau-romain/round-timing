package csrftoken

import "context"

type csrfTokenKey struct{}

// ContextKey is the key used to store the CSRF token in the context.
// Exported so the middleware can set it.
var ContextKey = csrfTokenKey{}

// FromContext retrieves the CSRF token from the context.
func FromContext(ctx context.Context) string {
	token, _ := ctx.Value(ContextKey).(string)
	return token
}
