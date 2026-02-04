package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/config"
	handlersPage "github.com/rousseau-romain/round-timing/handlers/page"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerPageRoutes(r *mux.Router, handler *handlersPage.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/", middleware.AllowToBeAuth(handler.HandleHome, authService, logger)).Methods("GET")
	r.Handle("/version", middleware.AllowToBeAuth(handler.HandleVersion, authService, logger)).Methods("GET")
	r.Handle("/privacy", middleware.AllowToBeAuth(handler.HandlePrivacy, authService, logger)).Methods("GET")
	r.Handle("/cgu", middleware.AllowToBeAuth(handler.HandleCGU, authService, logger)).Methods("GET")

	if config.ENV == "development" || config.ENV == "staging" {
		r.Handle("/test/ui", middleware.AllowToBeAuth(handler.HandleTestUI, authService, logger)).Methods("GET")
	}
}
