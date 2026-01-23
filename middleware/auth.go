package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	"github.com/rousseau-romain/round-timing/model/system"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/pkg/constants"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/views/components/layout"
	"github.com/rousseau-romain/round-timing/views/page"
)

// enabledUserIfWhiteListed checks if user should be enabled based on whitelist.
func enabledUserIfWhiteListed(w http.ResponseWriter, logger *slog.Logger, user userModel.User) bool {
	if system.GetFeatureFlagIsEnabled("WHITE_LIST") && !user.Enabled {
		isWhiteListed, err := userModel.IsEmailWhiteListed(user.Email)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		if isWhiteListed {
			var t = true
			userModel.UpdateUser(user.Id, userModel.UserUpdate{Enabled: &t})
			return true
		}
		errorTitle := "You can't acces here"
		errorMessage := fmt.Sprintf("Ask to be add to whitelist at email %s", constants.MailContact)
		logger.Info("User is not white listed!", "userEmail", user.Email)
		w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s&errorMessages=%s", url.QueryEscape(errorTitle), errorMessage))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return false
	}
	return true
}

// AllowToBeAuth allows both authenticated and unauthenticated users.
func AllowToBeAuth(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	}
}

// RequireAuth requires user to be authenticated.
func RequireAuth(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authService.GetAuthenticateUserFromRequest(r, logger)
		if err != nil {
			logger.Error(err.Error())
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		logger = logger.With("userId", user.Id)

		if !enabledUserIfWhiteListed(w, logger, user) {
			return
		}

		handlerFunc(w, r)
	}
}

// RequireAuthAndAdmin requires user to be authenticated and an admin.
func RequireAuthAndAdmin(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authService.GetAuthenticateUserFromRequest(r, logger)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		logger = logger.With("userId", user.Id)

		if !enabledUserIfWhiteListed(w, logger, user) {
			return
		}

		if !user.IsAdmin {
			errorTitle := "You can't acces here"
			logger.Info("User is not Admin", "userId", user.Id)
			w.Header().Set("Location", fmt.Sprintf("/?errorTitle=%s", url.QueryEscape(errorTitle)))
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		handlerFunc(w, r)
	}
}

// RequireNotAuth requires user to NOT be authenticated.
func RequireNotAuth(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := authService.GetAuthenticateUserFromRequest(r, logger)
		if user.Id == 0 {
			handlerFunc(w, r)
			return
		}
		logger = logger.With("userId", user.Id)

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

// RequireAuthAndSpectateOfUserMatch requires user to be authenticated and a spectator of the match.
func RequireAuthAndSpectateOfUserMatch(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authService.GetAuthenticateUserFromRequest(r, logger)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		logger = logger.With("userId", user.Id)

		if !enabledUserIfWhiteListed(w, logger, user) {
			return
		}

		vars := mux.Vars(r)
		matchId, _ := strconv.Atoi(vars["idMatch"])

		_, err = matchModel.GetMatch(matchId)
		if err != nil {
			languages, _ := system.GetLanguages()
			errorMessage := i18n.T(r.Context(), "page.match.errors.match-not-found", i18n.M{"matchId": matchId})
			logger.Error("Match not found", "matchId", matchId)
			w.WriteHeader(http.StatusNotFound)
			page.NotFoundPage(errorMessage, []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
			return
		}

		userMatch, err := userModel.GetUserIdByMatch(matchId)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isUsersSpectateByIdUser, err := userModel.IsUsersSpectateByIdUser(userMatch.Id, user.IdShare)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isUsersSpectateByIdUser {
			errorMessage := i18n.T(r.Context(), "page.match.errors.match-unauthorized-spectator", i18n.M{"matchId": matchId})
			logger.Info("User is not spectator for match", "userId", user.Id, "userMatchId", userMatch.Id, "matchId", matchId)
			w.WriteHeader(http.StatusForbidden)
			languages, _ := system.GetLanguages()
			page.ForbidenPage(errorMessage, []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
			return
		}

		handlerFunc(w, r)
	}
}

// RequireAuthAndHisMatch requires user to be authenticated and the owner of the match.
func RequireAuthAndHisMatch(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authService.GetAuthenticateUserFromRequest(r, logger)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		logger = logger.With("userId", user.Id)

		if !enabledUserIfWhiteListed(w, logger, user) {
			return
		}

		vars := mux.Vars(r)
		matchId, _ := strconv.Atoi(vars["idMatch"])

		_, err = matchModel.GetMatch(matchId)
		if err != nil {
			languages, _ := system.GetLanguages()
			errorMessage := i18n.T(r.Context(), "page.match.errors.match-not-found", i18n.M{"matchId": matchId})
			logger.Error("Match not found", "matchId", matchId)
			w.WriteHeader(http.StatusNotFound)
			page.NotFoundPage(errorMessage, []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
			return
		}

		userMatch, err := userModel.GetUserIdByMatch(matchId)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userMatch.Id != user.Id {
			languages, _ := system.GetLanguages()
			errorMessage := i18n.T(r.Context(), "page.match.errors.match-unauthorized", i18n.M{"matchId": matchId})
			logger.Info("User is not the owner of the match", "userId", user.Id, "userMatchId", userMatch.Id)
			w.WriteHeader(http.StatusUnauthorized)
			page.ForbidenPage(errorMessage, []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
			return
		}

		handlerFunc(w, r)
	}
}

// RequireAuthAndHisAccount requires user to be authenticated and accessing their own account.
func RequireAuthAndHisAccount(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authService.GetAuthenticateUserFromRequest(r, logger)
		if err != nil {
			logger.Error(err.Error())
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		logger = logger.With("userId", user.Id)

		if !enabledUserIfWhiteListed(w, logger, user) {
			return
		}

		vars := mux.Vars(r)
		userId, _ := strconv.Atoi(vars["idUser"])

		if user.Id != userId {
			logger.Info("User is not the owner of the account", "userId", user.Id, "userId", userId)
			http.Error(w, fmt.Sprintf("User %v is not the owner of the account %v", user.Id, userId), http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r)
	}
}
