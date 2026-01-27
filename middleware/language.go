package middleware

import (
	"log/slog"
	"net/http"

	"github.com/invopop/ctxi18n"
	"github.com/rousseau-romain/round-timing/model/system"
	"github.com/rousseau-romain/round-timing/pkg/lang"
	"github.com/rousseau-romain/round-timing/service/auth"
)

// Language wraps an http.Handler to set the locale based on user preference or browser settings.
func Language(handler http.Handler, authService *auth.AuthService, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := authService.GetAuthenticateUserFromRequest(r, logger)

		if user.Id != 0 {
			r = auth.RequestWithUser(r, user)
		}

		var locale string
		var err error

		if user.Id == 0 {
			locale = lang.GetPreferred(r)
		} else {
			locale, err = system.GetLanguageLocaleById(user.IdLanguage)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		ctx, err := ctxi18n.WithLocale(r.Context(), locale)
		if err != nil {
			logger.Error("error setting locale", "error", err)
			http.Error(w, "error setting locale", http.StatusBadRequest)
			return
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
