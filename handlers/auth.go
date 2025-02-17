package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rousseau-romain/round-timing/helper"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	page.SigninPage(GetPageNavDefault(r), h.languages, h.error).Render(r.Context(), w)
}

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Printf("User already authenticated! %v", u)

		page.SigninPage(GetPageNavDefault(r), h.languages, h.error).Render(r.Context(), w)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)

		return
	}

	userAlreadyExists, err := model.UserExists(user.UserID)

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if !userAlreadyExists {
		lang := helper.GetPreferredLanguage(r)

		id, err := model.GetLanguagesIdByCode(lang)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		_, err = model.CreateUser(model.UserCreate{
			Oauth2Id:   user.UserID,
			Email:      user.Email,
			IdLanguage: id,
		})
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	}

	err = h.auth.StoreUserSession(w, r, user)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	h.auth.RemoveUserSession(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
