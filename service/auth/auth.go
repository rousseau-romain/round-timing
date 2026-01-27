package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/model/user"
)

type contextKey string

const userContextKey contextKey = "authenticatedUser"

func RequestWithUser(r *http.Request, u user.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userContextKey, u))
}

func UserFromRequest(r *http.Request) (user.User, bool) {
	u, ok := r.Context().Value(userContextKey).(user.User)
	return u, ok && u.Id != 0
}

type AuthService struct{}

func NewAuthService(store sessions.Store) *AuthService {
	gothic.Store = store

	goth.UseProviders(
		google.New(
			config.GOOGLE_CLIENT_ID,
			config.GOOGLE_CLIENT_SECRET,
			buildCallbackURL("google"),
			"email",
		),
		discord.New(
			config.DISCORD_CLIENT_ID,
			config.DISCORD_CLIENT_SECRET,
			buildCallbackURL("discord"),
			"email",
		),
	)

	return &AuthService{}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (s *AuthService) StoreUserSession(w http.ResponseWriter, r *http.Request, slog *slog.Logger, user goth.User) error {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := gothic.Store.Get(r, SessionName)

	session.Values["user"] = user

	err := session.Save(r, w)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func (s *AuthService) RemoveUserSession(w http.ResponseWriter, r *http.Request, slog *slog.Logger) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,            // Ensure Secure flag is set if using HTTPS
		Expires:  time.Unix(0, 0), // Expire the cookie immediately
	})

	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = goth.User{}
	session.Options.MaxAge = -1

	session.Save(r, w)
}

var errNotAuthenticated = errors.New("user is not authenticated")

func (s *AuthService) GetAuthenticateUserFromRequest(r *http.Request, slog *slog.Logger) (user.User, error) {
	// Return cached user if already resolved for this request
	if u, ok := r.Context().Value(userContextKey).(user.User); ok && u.Id != 0 {
		return u, nil
	}

	// Try OAuth2 session first
	if u, err := s.getUserFromOAuthSession(r, slog); err == nil {
		return u, nil
	}

	// Fall back to JWT cookie (email/password login)
	if u, err := s.getUserFromJWT(r, slog); err == nil {
		return u, nil
	}

	return user.User{}, errNotAuthenticated
}

func (s *AuthService) getUserFromOAuthSession(r *http.Request, slog *slog.Logger) (user.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return user.User{}, err
	}

	sessionUser, ok := session.Values["user"]
	if !ok || sessionUser == nil {
		return user.User{}, errNotAuthenticated
	}

	gothUser, ok := sessionUser.(goth.User)
	if !ok || gothUser.UserID == "" {
		return user.User{}, errNotAuthenticated
	}

	u, err := user.GetUserByOauth2Id(gothUser.UserID)
	if err != nil {
		slog.Error("Error fetching user", "goticUserId", gothUser.UserID, "error", err)
		return user.User{}, err
	}
	return u, nil
}

func (s *AuthService) getUserFromJWT(r *http.Request, slog *slog.Logger) (user.User, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return user.User{}, errNotAuthenticated
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		slog.Error("Unauthorized: Invalid token")
		return user.User{}, errNotAuthenticated
	}

	u, err := user.GetUserByEmail(claims.Email)
	if err != nil {
		slog.Error("Error fetching user by email", "email", claims.Email, "error", err)
		return user.User{}, err
	}
	return u, nil
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", config.PUBLIC_HOST_PORT, provider)
}
