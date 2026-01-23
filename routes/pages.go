package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerPageRoutes(r *mux.Router, handler *handlers.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/", middleware.AllowToBeAuth(handler.HandlersHome, authService, logger)).Methods("GET")
	r.Handle("/commit-id", middleware.AllowToBeAuth(handler.HandlerCommitId, authService, logger)).Methods("GET")
	r.Handle("/version", middleware.AllowToBeAuth(handler.HandlerVersion, authService, logger)).Methods("GET")
	r.Handle("/privacy", middleware.AllowToBeAuth(handler.HandlerPrivacy, authService, logger)).Methods("GET")
	r.Handle("/cgu", middleware.AllowToBeAuth(handler.HandlerCGU, authService, logger)).Methods("GET")
}
