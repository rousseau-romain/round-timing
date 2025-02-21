package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/helper"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	page.SigninPage(GetPageNavDefault(r), h.languages, r.URL.Path, h.error).Render(r.Context(), w)
}

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Printf("User already authenticated! %v", u)

		page.SigninPage(GetPageNavDefault(r), h.languages, r.URL.Path, h.error).Render(r.Context(), w)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		log.Println(err)
		return
	}

	userAlreadyExists, err := model.UserExistsByOauth2Id(user.UserID)

	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	if !userAlreadyExists {
		providerLoginName, err := model.UserExistsByEmail(user.Email)

		if err != nil {
			fmt.Fprintln(w, err)
			log.Println(err)
			return
		}

		if providerLoginName != "" {
			errorTitle := i18n.T(r.Context(), "global.error") + " " + i18n.T(r.Context(), "page.signin.title")
			errorMessage := i18n.T(r.Context(), "page.signin.already-exists-with-provider", i18n.M{"email": user.Email, "provider": providerLoginName})
			err := gothic.Logout(w, r)
			if err != nil {
				log.Println(err)
				return
			}

			h.auth.RemoveUserSession(w, r)

			w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s&errorMessages=%s", url.QueryEscape(errorTitle), url.QueryEscape(errorMessage)))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		lang := helper.GetPreferredLanguage(r)

		id, err := model.GetLanguagesIdByCode(lang)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		_, err = model.CreateUser(model.UserCreate{
			ProviderLogin: user.Provider,
			Oauth2Id:      user.UserID,
			Email:         user.Email,
			IdLanguage:    id,
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
