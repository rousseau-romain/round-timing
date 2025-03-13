package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/helper"
	"github.com/rousseau-romain/round-timing/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
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
	token, _ := helper.GenerateSalt()
	csrfTokens[sessionID] = token
	return token
}

var csrfTokens = make(map[string]string) // Store CSRF tokens

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func enabledUserIfWhiteListed(w http.ResponseWriter, slog *slog.Logger, user model.User) bool {
	if model.GetFeatureFlagIsEnabled("WHITE_LIST") && !user.Enabled {
		isWhiteListed, err := model.IsEmailWhiteListed(user.Email)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		if isWhiteListed {
			var t = true
			model.UpdateUser(user.Id, model.UserUpdate{Enabled: &t})
			return true
		}
		errorTitle := "You can't acces here"
		errorMessage := fmt.Sprintf("Ask to be add to whitelist at email %s", helper.MailContact)
		slog.Info("User is not white listed!", "userEmail", user.Email)
		w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s&errorMessages=%s", url.QueryEscape(errorTitle), errorMessage))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return false
	}
	return true
}

func (s *AuthService) GetSessionUser(r *http.Request, slog *slog.Logger) (goth.User, error) {

	user, err := s.GetAuthenticateUserFromRequest(r, slog)

	return goth.User{
		Email:  user.Email,
		UserID: strconv.Itoa(user.Id),
	}, err

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

func AllowToBeAuth(handlerFunc http.HandlerFunc, auth *AuthService, slog *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
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

func (s *AuthService) GetAuthenticateUserFromRequest(r *http.Request, slog *slog.Logger) (model.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	var user model.User
	if err != nil {
		return user, err
	}

	u := session.Values["user"]

	// User is not from OAuth2 verify if user is authenticated by email
	if u == nil {
		cookieToken, err := getCookieHandler(r, "token")
		if err != nil {
			return user, errors.New("user is not authenticated")
		}

		claims := Claims{}
		token, err := jwt.ParseWithClaims(cookieToken.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET_KEY), nil
		})
		if err != nil || !token.Valid {
			slog.Error("Unauthorized: Invalid token", "token", cookieToken.Value)
			return user, err
		}

		// Extract claims and use them
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			user, err := model.GetUserByEmail(claims.Email)
			if err != nil {
				slog.Error("Error fetching user by email", "email", claims.Email, "error", err)
				return user, err
			}
			return user, nil
		}
		slog.Info("user is not authenticated", "userId", user.Id)
		return user, errors.New("user is not authenticated")
	}

	goticUser := u.(goth.User)
	userDb, err := model.GetUserByOauth2Id(goticUser.UserID)
	if err != nil {
		slog.Error("Error fetching user", "goticUserId", goticUser.UserID, "error", err)
		return userDb, err
	}
	return userDb, nil
}

func RequireAuth(handlerFunc http.HandlerFunc, auth *AuthService, slog *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetAuthenticateUserFromRequest(r, slog)
		if err != nil {
			slog.Error(err.Error())
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		if !enabledUserIfWhiteListed(w, slog, user) {
			return
		}

		handlerFunc(w, r)
	}
}

func RequireAuthAndAdmin(handlerFunc http.HandlerFunc, auth *AuthService, slog *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetAuthenticateUserFromRequest(r, slog)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		if !enabledUserIfWhiteListed(w, slog, user) {
			return
		}

		if !user.IsAdmin {
			errorTitle := "You can't acces here"
			slog.Error("User is not Admin", "userId", user.Id)
			w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s", url.QueryEscape(errorTitle)))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		handlerFunc(w, r)
	}
}

func RequireNotAuth(handlerFunc http.HandlerFunc, auth *AuthService, slog *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.GetAuthenticateUserFromRequest(r, slog)
		if user.Id == 0 {
			handlerFunc(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func RequireAuthAndSpectateOfUserMatch(handlerFunc http.HandlerFunc, auth *AuthService, slog *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetAuthenticateUserFromRequest(r, slog)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		if !enabledUserIfWhiteListed(w, slog, user) {
			return
		}

		vars := mux.Vars(r)

		matchId, _ := strconv.Atoi(vars["idMatch"])

		userMatch, err := model.GetUserIdByMatch(matchId)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isUsersSpectateByIdUser, err := model.IsUsersSpectateByIdUser(userMatch.Id, user.IdShare)

		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isUsersSpectateByIdUser {
			errorTitle := "You are not spectater"
			errorMessage := "You are not spectator of this match ask user to add you in list spectater in this profile"
			slog.Info("User is not spectater for match", "user", user.Id, "match", matchId)
			w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s&errorMessages=%s", url.QueryEscape(errorTitle), url.QueryEscape(errorMessage)))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		handlerFunc(w, r)
	}
}

func RequireAuthAndHisMatch(handlerFunc http.HandlerFunc, auth *AuthService, slog *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetAuthenticateUserFromRequest(r, slog)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		if !enabledUserIfWhiteListed(w, slog, user) {
			return
		}

		vars := mux.Vars(r)

		matchId, _ := strconv.Atoi(vars["idMatch"])

		userMatch, err := model.GetUserIdByMatch(matchId)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userMatch.Id != user.Id {
			slog.Error("User is not the owner of the match", "userId", user.Id, "userMatchId", userMatch.Id)
			http.Error(w, fmt.Sprintf("User %v is not the owner of the match %v", user.Id, userMatch.Id), http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r)
	}
}

func RequireAuthAndHisAccount(handlerFunc http.HandlerFunc, auth *AuthService, slog *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetAuthenticateUserFromRequest(r, slog)
		if err != nil {
			slog.Error(err.Error())
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}

		if !enabledUserIfWhiteListed(w, slog, user) {
			return
		}

		vars := mux.Vars(r)

		userId, _ := strconv.Atoi(vars["idUser"])

		if user.Id != userId {
			slog.Error("User is not the owner of the account", "userId", user.Id, "userId", userId)
			http.Error(w, fmt.Sprintf("User %v is not the owner of the account %v", user.Id, userId), http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r)
	}
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", config.PUBLIC_HOST_PORT, provider)
}
