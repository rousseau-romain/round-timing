package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	handlersTournament "github.com/rousseau-romain/round-timing/handlers/tournament"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerTournamentRoutes(r *mux.Router, handler *handlersTournament.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/tournament", middleware.RequireAuth(handler.HandleListTournament, authService, logger)).Methods("GET")
	r.Handle("/tournament", middleware.RequireAuth(handler.HandleCreateTournament, authService, logger)).Methods("POST")
	r.Handle("/tournament/{idTournament:[0-9]+}", middleware.RequireAuthAndHisTournament(handler.HandleViewTournament, authService, logger)).Methods("GET")
	r.Handle("/tournament/{idTournament:[0-9]+}", middleware.RequireAuthAndHisTournament(handler.HandleDeleteTournament, authService, logger)).Methods("DELETE")
	r.Handle("/tournament/archived", middleware.RequireAuth(handler.HandleListArchivedTournament, authService, logger)).Methods("GET")
	r.Handle("/tournament/archived/hide", middleware.RequireAuth(handler.HandleHideArchivedSection, authService, logger)).Methods("GET")
	r.Handle("/tournament/{idTournament:[0-9]+}/archive", middleware.RequireAuthAndHisTournament(handler.HandleArchiveTournament, authService, logger)).Methods("PATCH")
	r.Handle("/tournament/{idTournament:[0-9]+}/unarchive", middleware.RequireAuthAndHisTournament(handler.HandleUnarchiveTournament, authService, logger)).Methods("PATCH")

	// Teams
	r.Handle("/tournament/team", middleware.RequireAuth(handler.HandleAddTeam, authService, logger)).Methods("POST")
	r.Handle("/tournament/team/{idTeam:[0-9]+}", middleware.RequireAuth(handler.HandleDeleteTeam, authService, logger)).Methods("DELETE")

	// Players
	r.Handle("/tournament/player", middleware.RequireAuth(handler.HandleAddPlayer, authService, logger)).Methods("POST")
	r.Handle("/tournament/player/{idPlayer:[0-9]+}", middleware.RequireAuth(handler.HandleDeletePlayer, authService, logger)).Methods("DELETE")

	// Team Players (composition)
	r.Handle("/tournament/team/{idTeam:[0-9]+}/player", middleware.RequireAuth(handler.HandleAddTeamPlayer, authService, logger)).Methods("POST")
	r.Handle("/tournament/team/{idTeam:[0-9]+}/player/{idTeamPlayer:[0-9]+}", middleware.RequireAuth(handler.HandleDeleteTeamPlayer, authService, logger)).Methods("DELETE")

	// Matches
	r.Handle("/tournament/{idTournament:[0-9]+}/matchs", middleware.RequireAuthAndHisTournament(handler.HandleViewTournamentMatchs, authService, logger)).Methods("GET")
	r.Handle("/tournament/{idTournament:[0-9]+}/match", middleware.RequireAuthAndHisTournament(handler.HandleCreateMatch, authService, logger)).Methods("POST")
	r.Handle("/tournament/{idTournament:[0-9]+}/match/{idMatch:[0-9]+}/score", middleware.RequireAuthAndHisTournament(handler.HandleUpdateMatchScore, authService, logger)).Methods("POST")
	r.Handle("/tournament/{idTournament:[0-9]+}/match/{idMatch:[0-9]+}/score", middleware.RequireAuthAndHisTournament(handler.HandleIncrementMatchScore, authService, logger)).Methods("PATCH")
	r.Handle("/tournament/{idTournament:[0-9]+}/match/{idMatch:[0-9]+}/bo", middleware.RequireAuthAndHisTournament(handler.HandleEditMatchBo, authService, logger)).Methods("PATCH")
	r.Handle("/tournament/{idTournament:[0-9]+}/match/{idMatch:[0-9]+}/kills", middleware.RequireAuthAndHisTournament(handler.HandleUpdateMatchKills, authService, logger)).Methods("PATCH")
	r.Handle("/tournament/{idTournament:[0-9]+}/match/{idMatch:[0-9]+}/status", middleware.RequireAuthAndHisTournament(handler.HandleUpdateMatchStatus, authService, logger)).Methods("PATCH")
	r.Handle("/tournament/{idTournament:[0-9]+}/match/{idMatch:[0-9]+}", middleware.RequireAuthAndHisTournament(handler.HandleDeleteMatch, authService, logger)).Methods("DELETE")
	r.Handle("/tournament/{idTournament:[0-9]+}/matches-section", middleware.RequireAuthAndHisTournament(handler.HandleMatchesSection, authService, logger)).Methods("GET")
	r.Handle("/tournament/{idTournament:[0-9]+}/status-badge", middleware.RequireAuthAndHisTournament(handler.HandleTournamentStatusBadge, authService, logger)).Methods("GET")
	r.Handle("/tournament/{idTournament:[0-9]+}/generate-next-round", middleware.RequireAuthAndHisTournament(handler.HandleGenerateNextRound, authService, logger)).Methods("POST")
	r.Handle("/tournament/{idTournament:[0-9]+}/generate-next-round-ordered", middleware.RequireAuthAndHisTournament(handler.HandleGenerateNextRoundOrdered, authService, logger)).Methods("POST")
}
