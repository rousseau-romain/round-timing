package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerMatchRoutes(r *mux.Router, handler *handlers.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/match", middleware.RequireAuth(handler.HandlersListMatch, authService, logger)).Methods("GET")
	r.Handle("/match", middleware.RequireAuth(handler.HandlersCreateMatch, authService, logger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersDeleteMatch, authService, logger)).Methods("DELETE")
	r.Handle("/match/{idMatch:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersMatch, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/spectate", middleware.RequireAuthAndSpectateOfUserMatch(handler.HandlerSpectateMatch, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/start", middleware.RequireAuthAndHisMatch(handler.HandlerStartMatchPage, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/reset", middleware.RequireAuthAndHisMatch(handler.HandlerResetMatchPage, authService, logger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/increase-round", middleware.RequireAuthAndHisMatch(handler.HandlerMatchNextRound, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/table-live", middleware.AllowToBeAuth(handler.HandlerMatchTableLive, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/toggle-mastery/{toggleBool:[0-1]}", middleware.RequireAuthAndHisMatch(handler.HandlerToggleMatchMastery, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/use", middleware.RequireAuthAndHisMatch(handler.HandlerUsePlayerSpell, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/remove-round-recovery", middleware.RequireAuthAndHisMatch(handler.HandlerRemoveRoundRecoveryPlayerSpell, authService, logger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player", middleware.RequireAuthAndHisMatch(handler.HandlersCreatePlayer, authService, logger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersUpdatePlayer, authService, logger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersDeletePlayer, authService, logger)).Methods("DELETE")
}
