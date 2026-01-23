package auth

import (
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
	"github.com/rousseau-romain/round-timing/pkg/password"
)

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

func GenerateCSRFToken(sessionID string) string {
	token, _ := password.GenerateSalt()
	csrfTokens[sessionID] = token
	return token
}

var csrfTokens = make(map[string]string) // Store CSRF tokens

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

	// Expire the CSRF token cookie as well
	http.SetCookie(w, &http.Cookie{
		Name:    "csrf_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
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

func getCookieHandler(r *http.Request, name string) (http.Cookie, error) {
	// Retrieve a specific cookie by name
	cookie, err := r.Cookie(name)
	if err != nil {
		// Handle error if the cookie is not found
		if err == http.ErrNoCookie {
			return http.Cookie{}, errors.New("cookie not found")
		} else {
			return http.Cookie{}, fmt.Errorf("error retrieving cookie: %v", err)
		}
	}
	return *cookie, err
}

func (s *AuthService) GetAuthenticateUserFromRequest(r *http.Request, slog *slog.Logger) (user.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	var u user.User
	if err != nil {
		slog.Error(err.Error())
		return u, err
	}

	slog = slog.With("userEmail", u)

	sessionUser := session.Values["user"]

	// User is not from OAuth2 verify if user is authenticated by email
	if sessionUser == nil {
		cookieToken, err := getCookieHandler(r, "token")
		if err != nil {
			return u, errors.New("user is not authenticated")
		}

		claims := Claims{}
		token, err := jwt.ParseWithClaims(cookieToken.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET_KEY), nil
		})
		if err != nil || !token.Valid {
			slog.Error("Unauthorized: Invalid token", "token", cookieToken.Value)
			return u, err
		}

		// Extract claims and use them
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			foundUser, err := user.GetUserByEmail(claims.Email)
			if err != nil {
				slog.Error("Error fetching user by email", "email", claims.Email, "error", err)
				return foundUser, err
			}
			return foundUser, nil
		}
		slog.Info("user is not authenticated", "userId", u.Id)
		return u, errors.New("user is not authenticated")
	}

	goticUser := sessionUser.(goth.User)
	userDb, err := user.GetUserByOauth2Id(goticUser.UserID)
	if err != nil {
		slog.Error("Error fetching user", "goticUserId", goticUser.UserID, "error", err)
		return userDb, err
	}
	return userDb, nil
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", config.PUBLIC_HOST_PORT, provider)
}
