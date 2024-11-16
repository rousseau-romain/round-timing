package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"round-timing/config"
	"round-timing/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
)

type AuthService struct{}

func NewAuthService(store sessions.Store) *AuthService {
	gothic.Store = store

	goth.UseProviders(
		github.New(
			config.GITHUB_CLIENT_ID,
			config.GITHUB_CLIENT_SECRET,
			buildCallbackURL("github"),
		),
		discord.New(
			config.DISCORD_CLIENT_ID,
			config.DISCORD_CLIENT_SECRET,
			buildCallbackURL("discord"),
		),
	)

	return &AuthService{}
}

func (s *AuthService) GetSessionUser(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return goth.User{}, err
	}

	u := session.Values["user"]
	if u == nil {
		return goth.User{}, fmt.Errorf("user is not authenticated! %v", u)
	}

	return u.(goth.User), nil
}

func (s *AuthService) StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := gothic.Store.Get(r, SessionName)

	session.Values["user"] = user

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func (s *AuthService) RemoveUserSession(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = goth.User{}
	// delete the cookie immediately
	session.Options.MaxAge = -1

	session.Save(r, w)
}

func AllowToBeAuth(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
}

func RequireAuth(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		user, err := model.GetUserByOauth2Id(session.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !user.Enabled {
			log.Println("User is not enabled!")
			http.Redirect(w, r, "?require-user-enabled", http.StatusTemporaryRedirect)
			return
		}

		log.Printf("user is authenticated! user: %v!", session.FirstName)

		handlerFunc(w, r)
	}
}

func RequireNotAuth(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			handlerFunc(w, r)
			return
		}

		log.Printf("user is authenticated! user: %v!", session.FirstName)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func RequireAuthAndHisMatch(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		user, err := model.GetUserByOauth2Id(session.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !user.Enabled {
			log.Println("User is not enabled!")
			http.Redirect(w, r, "?require-user-enabled", http.StatusTemporaryRedirect)
			return
		}

		vars := mux.Vars(r)

		matchId, _ := strconv.Atoi(vars["idMatch"])

		userMatch, err := model.GetUserIdByMatch(matchId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userMatch.Id != user.Id {
			log.Printf("User %v is not the owner of the match %v", user.Id, userMatch.Id)
			http.Error(w, fmt.Sprintf("User %v is not the owner of the match %v", user.Id, userMatch.Id), http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r)
	}
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", config.PUBLIC_HOST_PORT, provider)
}
