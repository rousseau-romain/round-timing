package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerPageRoutes(r *mux.Router, handler *handlers.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/", middleware.AllowToBeAuth(handler.HandleHome, authService, logger)).Methods("GET")
	r.Handle("/commit-id", middleware.AllowToBeAuth(handler.HandleCommitId, authService, logger)).Methods("GET")
	r.Handle("/version", middleware.AllowToBeAuth(handler.HandleVersion, authService, logger)).Methods("GET")
	r.Handle("/privacy", middleware.AllowToBeAuth(handler.HandlePrivacy, authService, logger)).Methods("GET")
	r.Handle("/cgu", middleware.AllowToBeAuth(handler.HandleCGU, authService, logger)).Methods("GET")
}
