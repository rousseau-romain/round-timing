package handlers

import (
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/markbates/goth/gothic"
	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/model/system"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/pkg/lang"
	"github.com/rousseau-romain/round-timing/pkg/password"
	"github.com/rousseau-romain/round-timing/service/auth"
	authPage "github.com/rousseau-romain/round-timing/views/page/auth"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	authPage.SigninPage(GetPageNavDefault(r), h.languages, r.URL.Path, h.error).Render(r.Context(), w)
}

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		h.Slog.Info("User already authenticated", "goticUserId", u.UserID)
		authPage.SigninPage(GetPageNavDefault(r), h.languages, r.URL.Path, h.error).Render(r.Context(), w)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		h.Slog.Error(err.Error())
		return
	}

	userAlreadyExists, err := userModel.UserExistsByOauth2Id(user.UserID)

	if err != nil {
		h.Slog.Error(err.Error())
		fmt.Fprintln(w, err)
		return
	}

	h.Slog = h.Slog.With("userOauth2Id", user.UserID)

	if !userAlreadyExists {
		providerLoginName, err := userModel.UserExistsByEmail(user.Email)

		if err != nil {
			fmt.Fprintln(w, err)
			h.Slog.Error(err.Error())
			return
		}

		if providerLoginName != "" {
			errorTitle := i18n.T(r.Context(), "global.error") + " " + i18n.T(r.Context(), "page.signin.title")
			errorMessage := i18n.T(r.Context(), "page.signin.already-exists-with-provider", i18n.M{"email": user.Email, "provider": providerLoginName})
			err := gothic.Logout(w, r)
			if err != nil {
				h.Slog.Error(err.Error())
				return
			}

			h.auth.RemoveUserSession(w, r, h.Slog)

			w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s&errorMessages=%s", url.QueryEscape(errorTitle), url.QueryEscape(errorMessage)))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		locale := lang.GetPreferred(r)

		idLanguage, err := system.GetLanguagesIdByCode(locale)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		_, err = userModel.CreateUser(userModel.UserCreate{
			ProviderLogin: user.Provider,
			Oauth2Id:      user.UserID,
			Email:         user.Email,
			IdLanguage:    idLanguage,
		})
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	}

	err = h.auth.StoreUserSession(w, r, h.Slog, user)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) HandleSignupEmail(w http.ResponseWriter, r *http.Request) {
	authPage.SignupPage(GetPageNavDefault(r), h.languages, r.URL.Path, h.error).Render(r.Context(), w)
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	pwd := strings.TrimSpace(r.FormValue("password"))
	pwdConfirmation := strings.TrimSpace(r.FormValue("password-confirmation"))

	var errMessages = []string{}
	if _, err := mail.ParseAddress(email); err != nil {
		errMessages = append(errMessages, i18n.T(r.Context(), "page.signup.error.email.not-valid"))
	}
	h.Slog = h.Slog.With("email", email)

	passwordIsValid, errorMessages := password.Validate(r, pwd)
	if !passwordIsValid {
		errMessages = append(errMessages, errorMessages...)
	}

	if pwd != pwdConfirmation {
		errMessages = append(errMessages, i18n.T(r.Context(), "page.signup.error.password.confirmation"))
	}

	if len(errMessages) > 0 {
		RenderComponentError(
			i18n.T(r.Context(), "page.signup.error.title"),
			errMessages,
			http.StatusBadRequest, w, r,
		)
		h.Slog.Info(strings.Join(errorMessages, "\n"))
		return
	}

	provider, err := userModel.UserExistsByEmail(email)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, "can't create user", http.StatusInternalServerError)
		return
	}
	if provider != "" {
		RenderComponentError(
			i18n.T(r.Context(), "page.signup.error.email.already-exists", i18n.M{"email": email, "provider": provider}),
			[]string{""},
			http.StatusConflict, w, r,
		)
		h.Slog.Info("Email already exist", "email", email, "provider", provider)
		return
	}

	salt, err := password.GenerateSalt()
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, "can't create user", http.StatusInternalServerError)
		return
	}

	hashedPassword := password.Hash(pwd, salt)

	locale := lang.GetPreferred(r)
	idLanguage, err := system.GetLanguagesIdByCode(locale)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, "can't create user", http.StatusInternalServerError)
		return
	}

	user := userModel.UserCreate{
		ProviderLogin: "email",
		Email:         email,
		Hash:          hashedPassword,
		IdLanguage:    idLanguage,
	}

	_, err = userModel.CreateUser(user)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, "can't create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/signin")
	w.WriteHeader(http.StatusCreated)
}

func generateToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &auth.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT_SECRET_KEY))
}

func (h *Handler) HandleLoginEmail(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	pwd := strings.TrimSpace(r.FormValue("password"))

	user, err := userModel.GetUserByEmail(email)

	if err != nil || !password.Check(user.Hash, pwd) {
		errMessage := i18n.T(r.Context(), "page.signin.invalid-credentials")
		RenderComponentError(
			errMessage,
			[]string{errMessage},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Info("Invalid credentials: " + errMessage)
		return
	}
	h.Slog = h.Slog.With("userId", user.Id)

	token, err := generateToken(user.Email)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	csrfToken := auth.GenerateCSRFToken(user.Email)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "csrf_token",
		Value:   csrfToken,
		Path:    "/",
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)
	err := gothic.Logout(w, r)
	if err != nil {
		h.Slog.Error(err.Error())
		return
	}

	h.auth.RemoveUserSession(w, r, h.Slog)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
