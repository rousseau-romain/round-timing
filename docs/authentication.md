# Authentication

## Overview

RoundTiming uses a dual authentication system:

1. **OAuth2** - Social login via Discord and Google (using Goth library)
2. **Email/Password** - Traditional login with JWT tokens

## Authentication Flow

### OAuth2 Flow (Discord/Google)

```text
User clicks "Sign in with Discord/Google"
    │
    ▼
GET /auth/{provider}
    │
    ▼
Redirect to OAuth provider
    │
    ▼
User authenticates
    │
    ▼
GET /auth/{provider}/callback
    │
    ▼
Create/update user in database
    │
    ▼
Store user in session (gorilla/sessions)
    │
    ▼
Redirect to home page
```

### Email/Password Flow

```text
User submits login form
    │
    ▼
POST /signin
    │
    ▼
Verify credentials
    │
    ▼
Generate JWT token
    │
    ▼
Set token cookie
    │
    ▼
Redirect to home page
```

## Session Management

Sessions are managed using `gorilla/sessions` with cookie storage.

### Configuration

Environment variables in `.env`:

```text
COOKIES_AUTH_SECRET=<secret-key>
COOKIES_AUTH_AGE_IN_SECONDS=86400
COOKIES_AUTH_IS_SECURE=true
COOKIES_AUTH_IS_HTTP_ONLY=true
JWT_SECRET_KEY=<your-secret-key>
```

### Session Store

```go
// service/auth/session.go
store := sessions.NewCookieStore([]byte(opts.CookiesKey))
store.MaxAge(opts.MaxAge)
store.Options.Path = "/"
store.Options.HttpOnly = opts.HttpOnly
store.Options.Secure = opts.Secure
```

## OAuth2 Providers

### Discord

```text
DISCORD_CLIENT_ID=<your-client-id>
DISCORD_CLIENT_SECRET=<your-client-secret>
```

Callback URL: `{PUBLIC_HOST_PORT}/auth/discord/callback`

### Google

```text
GOOGLE_CLIENT_ID=<your-client-id>
GOOGLE_CLIENT_SECRET=<your-client-secret>
```

Callback URL: `{PUBLIC_HOST_PORT}/auth/google/callback`

## JWT Tokens

Used for email/password authentication.

### Structure

```go
type Claims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}
```

## Route Protection

### Middleware Functions

Located in `service/auth/auth.go`:

| Function | Description |
| -------- | ----------- |
| `AllowToBeAuth` | No authentication required |
| `RequireAuth` | User must be authenticated |
| `RequireNotAuth` | User must NOT be authenticated (login/signup pages) |
| `RequireAuthAndAdmin` | User must be admin |
| `RequireAuthAndHisMatch` | User must own the match |
| `RequireAuthAndSpectateOfUserMatch` | User must be a spectator |
| `RequireAuthAndHisAccount` | User must own the account |

### Usage Example

```go
// Public route
r.Handle("/", auth.AllowToBeAuth(handler.HandlersHome, authService, logger))

// Protected route
r.Handle("/match", auth.RequireAuth(handler.HandlersListMatch, authService, logger))

// Admin only
r.Handle("/admin/user", auth.RequireAuthAndAdmin(handler.HandlersListUser, authService, logger))
```

## Whitelist Feature

When the `WHITE_LIST` feature flag is enabled:

1. New users are created as disabled
2. Users must be in `email_white_listed` table to be enabled
3. Non-whitelisted users see an error message

## CSRF Protection

CSRF tokens are generated per session:

```go
func GenerateCSRFToken(sessionID string) string {
    token, _ := helper.GenerateSalt()
    csrfTokens[sessionID] = token
    return token
}
```

## Logout

Logout clears:

- JWT token cookie
- CSRF token cookie
- Session data

```go
// service/auth/auth.go
func (s *AuthService) RemoveUserSession(w http.ResponseWriter, r *http.Request, slog *slog.Logger)
```

## Getting Current User

To get the authenticated user in a handler:

```go
user, err := auth.GetAuthenticateUserFromRequest(r, slog)
if err != nil {
    // User is not authenticated
}
// Use user.Id, user.Email, user.IsAdmin, etc.
```
