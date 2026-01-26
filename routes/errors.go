package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerErrorRoutes(r *mux.Router, handler *handlers.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/404", middleware.AllowToBeAuth(handler.HandleNotFound, authService, logger)).Methods("GET")
	r.Handle("/403", middleware.AllowToBeAuth(handler.HandleForbidden, authService, logger)).Methods("GET")
}
