package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	handlersPage "github.com/rousseau-romain/round-timing/handlers/page"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerErrorRoutes(r *mux.Router, handler *handlersPage.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/404", middleware.AllowToBeAuth(handler.HandleNotFound, authService, logger)).Methods("GET")
	r.Handle("/403", middleware.AllowToBeAuth(handler.HandleForbidden, authService, logger)).Methods("GET")
}
