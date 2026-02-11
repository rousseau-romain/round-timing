package auth

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/markbates/goth/gothic"
	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/model/system"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/pkg/lang"
	"github.com/rousseau-romain/round-timing/pkg/password"
	serviceAuth "github.com/rousseau-romain/round-timing/service/auth"
	authPage "github.com/rousseau-romain/round-timing/views/page/auth"
)

const authFailedMessage = "Authentication failed"

type Handler struct {
	*handlers.Handler
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	authPage.SigninPage(handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, h.Error).Render(r.Context(), w)
}

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		h.Slog.Info("User already authenticated", "goticUserId", u.UserID)
		authPage.SigninPage(handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, h.Error).Render(r.Context(), w)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		handlers.RespondWithError(w, r, h.Slog, err, authFailedMessage, http.StatusInternalServerError)
		return
	}

	userAlreadyExists, err := userModel.UserExistsByOauth2Id(r.Context(), user.UserID)

	if err != nil {
		handlers.RespondWithError(w, r, h.Slog, err, authFailedMessage, http.StatusInternalServerError)
		return
	}

	logger := h.Slog.With("userOauth2Id", user.UserID)

	if !userAlreadyExists {
		providerLoginName, err := userModel.UserExistsByEmail(r.Context(), user.Email)

		if err != nil {
			handlers.RespondWithError(w, r, logger, err, authFailedMessage, http.StatusInternalServerError)
			return
		}

		if providerLoginName != "" {
			errorTitle := i18n.T(r.Context(), "global.error") + " " + i18n.T(r.Context(), "page.signin.title")
			errorMessage := i18n.T(r.Context(), "page.signin.already-exists-with-provider", i18n.M{"email": user.Email, "provider": providerLoginName})
			err := gothic.Logout(w, r)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			h.Auth.RemoveUserSession(w, r, h.Slog)

			serviceAuth.SetFlash(w, r, errorTitle, []string{errorMessage}, "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		locale := lang.GetPreferred(r)

		idLanguage, err := system.GetLanguagesIdByCode(r.Context(), locale)
		if err != nil {
			handlers.RespondWithError(w, r, logger, err, authFailedMessage, http.StatusInternalServerError)
			return
		}
		_, err = userModel.CreateUser(r.Context(), userModel.UserCreate{
			ProviderLogin: user.Provider,
			Oauth2Id:      user.UserID,
			Email:         user.Email,
			IdLanguage:    idLanguage,
		})
		if err != nil {
			handlers.RespondWithError(w, r, logger, err, authFailedMessage, http.StatusInternalServerError)
			return
		}

		logger.Info("user created", "email", user.Email, "provider", user.Provider)
	}

	err = h.Auth.StoreUserSession(w, r, h.Slog, user)
	if err != nil {
		handlers.RespondWithError(w, r, h.Slog, err, authFailedMessage, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) HandleSignupEmail(w http.ResponseWriter, r *http.Request) {
	authPage.SignupPage(handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, h.Error).Render(r.Context(), w)
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handlers.RespondWithError(w, r, h.Slog, err, "Invalid request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	pwd := strings.TrimSpace(r.FormValue("password"))
	pwdConfirmation := strings.TrimSpace(r.FormValue("password-confirmation"))

	var errMessages = []string{}
	if _, err := mail.ParseAddress(email); err != nil {
		errMessages = append(errMessages, i18n.T(r.Context(), "page.signup.error.email.not-valid"))
	}
	logger := h.Slog.With("email", email)

	passwordIsValid, errorMessages := password.Validate(r, pwd)
	if !passwordIsValid {
		errMessages = append(errMessages, errorMessages...)
	}

	if pwd != pwdConfirmation {
		errMessages = append(errMessages, i18n.T(r.Context(), "page.signup.error.password.confirmation"))
	}

	if len(errMessages) > 0 {
		handlers.RenderComponentError(
			i18n.T(r.Context(), "page.signup.error.title"),
			errMessages,
			http.StatusBadRequest, w, r,
		)
		logger.Info(strings.Join(errorMessages, "\n"))
		return
	}

	provider, err := userModel.UserExistsByEmail(r.Context(), email)
	if err != nil {
		logger.Error(err.Error())
		httpError.CantCreateUser(w)
		return
	}
	if provider != "" {
		handlers.RenderComponentError(
			i18n.T(r.Context(), "page.signup.error.email.already-exists", i18n.M{"email": email, "provider": provider}),
			[]string{""},
			http.StatusConflict, w, r,
		)
		logger.Info("Email already exist", "email", email, "provider", provider)
		return
	}

	salt, err := password.GenerateSalt()
	if err != nil {
		logger.Error(err.Error())
		httpError.CantCreateUser(w)
		return
	}

	hashedPassword := password.Hash(pwd, salt)

	locale := lang.GetPreferred(r)
	idLanguage, err := system.GetLanguagesIdByCode(r.Context(), locale)
	if err != nil {
		logger.Error(err.Error())
		httpError.CantCreateUser(w)
		return
	}

	user := userModel.UserCreate{
		ProviderLogin: "email",
		Email:         email,
		Hash:          hashedPassword,
		IdLanguage:    idLanguage,
	}

	_, err = userModel.CreateUser(r.Context(), user)
	if err != nil {
		logger.Error(err.Error())
		httpError.CantCreateUser(w)
		return
	}

	logger.Info("user created", "email", email)

	w.Header().Set("HX-Redirect", "/signin")
	w.WriteHeader(http.StatusCreated)
}

func generateToken(email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.COOKIES_AUTH_AGE_IN_SECONDS) * time.Second)
	claims := &serviceAuth.Claims{
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
		handlers.RespondWithError(w, r, h.Slog, err, "Invalid request", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	pwd := strings.TrimSpace(r.FormValue("password"))

	user, err := userModel.GetUserByEmail(r.Context(), email)

	// Always run the hash check to prevent timing-based email enumeration.
	// When the user is not found, check against a dummy hash so the
	// response time is the same as for a wrong password.
	hashToCheck := user.Hash
	if err != nil {
		hashToCheck = password.DummyHash
	}
	if err != nil || !password.Check(hashToCheck, pwd) {
		errMessage := i18n.T(r.Context(), "page.signin.invalid-credentials")
		handlers.RenderComponentError(
			errMessage,
			[]string{errMessage},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Info("Invalid credentials: " + errMessage)
		return
	}
	logger := h.Slog.With("userId", user.Id)

	token, err := generateToken(user.Email)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   config.COOKIES_AUTH_AGE_IN_SECONDS,
	})

	logger.Info("user logged in", "email", email)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	user, _ := serviceAuth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	err := gothic.Logout(w, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	h.Auth.RemoveUserSession(w, r, h.Slog)

	logger.Info("user logged out")

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
