package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

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

func enabledUserIfWhiteListed(w http.ResponseWriter, user model.User) {
	if model.GetFeatureFlagIsEnabled("WHITE_LIST") && !user.Enabled {
		isWhiteListed, err := model.IsEmailWhiteListed(user.Email)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if isWhiteListed {
			var t = true
			model.UpdateUser(user.Id, model.UserUpdate{Enabled: &t})
			return
		}
		errorTitle := "You can't acces here"
		errorMessage := fmt.Sprintf("Ask to be add to whitelist at email %s", helper.MailContact)
		log.Printf("User (%v) is not white listed!", user.Email)
		w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s&errorMessage=%s", url.QueryEscape(errorTitle), errorMessage))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
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
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enabledUserIfWhiteListed(w, user)

		handlerFunc(w, r)
	}
}

func RequireAuthAndAdmin(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		user, err := model.GetUserByOauth2Id(session.UserID)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enabledUserIfWhiteListed(w, user)

		if !user.IsAdmin {
			errorTitle := "You can't acces here"
			log.Printf("%s", fmt.Sprintf("User (%d) is not Admin", user.Id))
			w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s", url.QueryEscape(errorTitle)))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		handlerFunc(w, r)
	}
}

func RequireNotAuth(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.GetSessionUser(r)
		if err != nil {
			handlerFunc(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func RequireAuthAndSpectateOfUserMatch(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		user, err := model.GetUserByOauth2Id(session.UserID)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enabledUserIfWhiteListed(w, user)

		vars := mux.Vars(r)

		matchId, _ := strconv.Atoi(vars["idMatch"])

		userMatch, err := model.GetUserIdByMatch(matchId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isUsersSpectateByIdUser, err := model.IsUsersSpectateByIdUser(userMatch.Id, user.IdShare)

		log.Println(user.Id, userMatch.IdShare, isUsersSpectateByIdUser)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isUsersSpectateByIdUser {
			errorTitle := "You are not spectater"
			errorMessage := "You are not spectator of this match ask user to add you in list spectater in this profile"
			log.Printf("%s", fmt.Sprintf("User (%d) is not spectater for match (%d)", user.Id, matchId))
			w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s&errorMessages=%s", url.QueryEscape(errorTitle), url.QueryEscape(errorMessage)))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		handlerFunc(w, r)
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
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enabledUserIfWhiteListed(w, user)

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

func RequireAuthAndHisAccount(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		user, err := model.GetUserByOauth2Id(session.UserID)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enabledUserIfWhiteListed(w, user)

		vars := mux.Vars(r)

		userId, _ := strconv.Atoi(vars["idUser"])

		if user.Id != userId {
			log.Printf("User %v is not the owner of the account %v", user.Id, userId)
			http.Error(w, fmt.Sprintf("User %v is not the owner of the account %v", user.Id, userId), http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r)
	}
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", config.PUBLIC_HOST_PORT, provider)
}
