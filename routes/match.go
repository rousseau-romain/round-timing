package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerMatchRoutes(r *mux.Router, handler *handlers.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/match", middleware.RequireAuth(handler.HandleListMatch, authService, logger)).Methods("GET")
	r.Handle("/match", middleware.RequireAuth(handler.HandleCreateMatch, authService, logger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandleDeleteMatch, authService, logger)).Methods("DELETE")
	r.Handle("/match/{idMatch:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandleMatch, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/spectate", middleware.RequireAuthAndSpectateOfUserMatch(handler.HandleSpectateMatch, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/start", middleware.RequireAuthAndHisMatch(handler.HandleStartMatch, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/reset", middleware.RequireAuthAndHisMatch(handler.HandleResetMatch, authService, logger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/increase-round", middleware.RequireAuthAndHisMatch(handler.HandleMatchNextRound, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/table-live", middleware.AllowToBeAuth(handler.HandleMatchTableLive, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/toggle-mastery/{toggleBool:[0-1]}", middleware.RequireAuthAndHisMatch(handler.HandleToggleMatchMastery, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/use", middleware.RequireAuthAndHisMatch(handler.HandleUsePlayerSpell, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/remove-round-recovery", middleware.RequireAuthAndHisMatch(handler.HandleRemoveRoundRecoveryPlayerSpell, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player", middleware.RequireAuthAndHisMatch(handler.HandleCreatePlayer, authService, logger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandleUpdatePlayer, authService, logger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandleDeletePlayer, authService, logger)).Methods("DELETE")
}
