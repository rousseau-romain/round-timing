package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	"github.com/rousseau-romain/round-timing/model/system"
	tournamentModel "github.com/rousseau-romain/round-timing/model/tournament"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/pkg/constants"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/views/components/layout"
	"github.com/rousseau-romain/round-timing/views/page"
)


// enabledUserIfWhiteListed checks if user should be enabled based on whitelist.
func enabledUserIfWhiteListed(ctx context.Context, w http.ResponseWriter, r *http.Request, logger *slog.Logger, user userModel.User) bool {
	if system.GetFeatureFlagIsEnabled(ctx, "WHITE_LIST") && !user.Enabled {
		isWhiteListed, err := userModel.IsEmailWhiteListed(ctx, user.Email)
		if err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return false
		}
		if isWhiteListed {
			var t = true
			userModel.UpdateUser(ctx, user.Id, userModel.UserUpdate{Enabled: &t})
			return true
		}
		errorTitle := "You can't acces here"
		errorMessage := fmt.Sprintf("Ask to be add to whitelist at email %s", constants.MailContact)
		logger.Info("User is not white listed!", "userEmail", user.Email)
		auth.SetFlash(w, r, errorTitle, []string{errorMessage}, "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

		if !enabledUserIfWhiteListed(r.Context(), w, r, logger, user) {
			return
		}

		handlerFunc(w, auth.RequestWithUser(r, user))
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

		if !enabledUserIfWhiteListed(r.Context(), w, r, logger, user) {
			return
		}

		if !user.IsAdmin {
			logger.Info("User is not Admin", "userId", user.Id)
			auth.SetFlash(w, r, "You can't acces here", []string{}, "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		handlerFunc(w, auth.RequestWithUser(r, user))
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

// getMatchAndOwner extracts matchId from URL vars, validates the match exists, and returns its owner.
func getMatchAndOwner(w http.ResponseWriter, r *http.Request, logger *slog.Logger, user userModel.User) (int, userModel.User, bool) {
	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	_, err := matchModel.GetMatch(r.Context(), matchId)
	if err != nil {
		languages, _ := system.GetLanguages(r.Context())
		errorMessage := i18n.T(r.Context(), "page.match.errors.match-not-found", i18n.M{"matchId": matchId})
		logger.Error("Match not found", "matchId", matchId)
		w.WriteHeader(http.StatusNotFound)
		page.NotFoundPage(errorMessage, []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
		return 0, userModel.User{}, false
	}

	userMatch, err := userModel.GetUserIdByMatch(r.Context(), matchId)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return 0, userModel.User{}, false
	}

	return matchId, userMatch, true
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

		if !enabledUserIfWhiteListed(r.Context(), w, r, logger, user) {
			return
		}

		matchId, userMatch, ok := getMatchAndOwner(w, r, logger, user)
		if !ok {
			return
		}

		isUsersSpectateByIdUser, err := userModel.IsUsersSpectateByIdUser(r.Context(), userMatch.Id, user.IdShare)
		if err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}

		if !isUsersSpectateByIdUser {
			errorMessage := i18n.T(r.Context(), "page.match.errors.match-unauthorized-spectator", i18n.M{"matchId": matchId})
			logger.Info("User is not spectator for match", "userId", user.Id, "userMatchId", userMatch.Id, "matchId", matchId)
			w.WriteHeader(http.StatusForbidden)
			languages, _ := system.GetLanguages(r.Context())
			page.ForbidenPage(errorMessage, []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
			return
		}

		handlerFunc(w, auth.RequestWithUser(r, user))
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

		if !enabledUserIfWhiteListed(r.Context(), w, r, logger, user) {
			return
		}

		matchId, userMatch, ok := getMatchAndOwner(w, r, logger, user)
		if !ok {
			return
		}

		if userMatch.Id != user.Id {
			languages, _ := system.GetLanguages(r.Context())
			errorMessage := i18n.T(r.Context(), "page.match.errors.match-unauthorized", i18n.M{"matchId": matchId})
			logger.Info("User is not the owner of the match", "userId", user.Id, "userMatchId", userMatch.Id)
			w.WriteHeader(http.StatusUnauthorized)
			page.ForbidenPage(errorMessage, []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
			return
		}

		handlerFunc(w, auth.RequestWithUser(r, user))
	}
}

// getTournamentAndOwner extracts tournamentId from URL vars, validates the tournament exists, and returns its owner ID.
func getTournamentAndOwner(w http.ResponseWriter, r *http.Request, logger *slog.Logger, user userModel.User) (int, int, bool) {
	vars := mux.Vars(r)
	tournamentId, _ := strconv.Atoi(vars["idTournament"])

	ownerId, err := tournamentModel.GetTournamentOwnerId(r.Context(), tournamentId)
	if err != nil {
		languages, _ := system.GetLanguages(r.Context())
		logger.Error("Tournament not found", "tournamentId", tournamentId)
		w.WriteHeader(http.StatusNotFound)
		page.NotFoundPage(fmt.Sprintf("Tournament id (%d) not found", tournamentId), []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
		return 0, 0, false
	}

	return tournamentId, ownerId, true
}

// RequireAuthAndHisTournament requires user to be authenticated and the owner of the tournament.
func RequireAuthAndHisTournament(handlerFunc http.HandlerFunc, authService *auth.AuthService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authService.GetAuthenticateUserFromRequest(r, logger)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
		logger = logger.With("userId", user.Id)

		if !enabledUserIfWhiteListed(r.Context(), w, r, logger, user) {
			return
		}

		tournamentId, ownerId, ok := getTournamentAndOwner(w, r, logger, user)
		if !ok {
			return
		}

		if ownerId != user.Id {
			languages, _ := system.GetLanguages(r.Context())
			logger.Info("User is not the owner of the tournament", "userId", user.Id, "tournamentId", tournamentId)
			w.WriteHeader(http.StatusUnauthorized)
			page.ForbidenPage(fmt.Sprintf("You don't have the permission to view this tournament (%d)", tournamentId), []layout.NavItem{}, languages, r.URL.Path, user).Render(r.Context(), w)
			return
		}

		handlerFunc(w, auth.RequestWithUser(r, user))
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

		if !enabledUserIfWhiteListed(r.Context(), w, r, logger, user) {
			return
		}

		vars := mux.Vars(r)
		userId, _ := strconv.Atoi(vars["idUser"])

		if user.Id != userId {
			logger.Info("User is not the owner of the account", "userId", user.Id, "targetUserId", userId)
			httpError.Forbidden(w)
			return
		}

		handlerFunc(w, auth.RequestWithUser(r, user))
	}
}
