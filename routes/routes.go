package routes

import (
	"net/http"

	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func Setup(handler *handlers.Handler, authService *auth.AuthService, logger *slog.Logger) *mux.Router {
	r := mux.NewRouter()

	registerPublicRoutes(r)
	registerPageRoutes(r, handler, authService, logger)
	registerMatchRoutes(r, handler, authService, logger)
	registerProfileRoutes(r, handler, authService, logger)
	registerAuthRoutes(r, handler, authService, logger)
	registerAdminRoutes(r, handler, authService, logger)
	registerErrorRoutes(r, handler, authService, logger)

	r.NotFoundHandler = http.HandlerFunc(handler.HandlersNotFound)

	return r
}
