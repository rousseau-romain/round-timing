package routes

import (
	"net/http"

	"log/slog"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	handlersAdmin "github.com/rousseau-romain/round-timing/handlers/admin"
	handlersAuth "github.com/rousseau-romain/round-timing/handlers/auth"
	handlersMatch "github.com/rousseau-romain/round-timing/handlers/match"
	handlersPage "github.com/rousseau-romain/round-timing/handlers/page"
	handlersProfile "github.com/rousseau-romain/round-timing/handlers/profile"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func Setup(handler *handlers.Handler, authService *auth.AuthService, logger *slog.Logger) *mux.Router {
	r := mux.NewRouter()

	matchH := &handlersMatch.Handler{Handler: handler}
	adminH := &handlersAdmin.Handler{Handler: handler}
	authH := &handlersAuth.Handler{Handler: handler}
	pageH := &handlersPage.Handler{Handler: handler}
	profileH := &handlersProfile.Handler{Handler: handler}

	registerPublicRoutes(r)
	registerPageRoutes(r, pageH, authService, logger)
	registerMatchRoutes(r, matchH, authService, logger)
	registerProfileRoutes(r, profileH, authService, logger)
	registerAuthRoutes(r, authH, authService, logger)
	registerAdminRoutes(r, adminH, authService, logger)
	registerErrorRoutes(r, pageH, authService, logger)

	r.NotFoundHandler = http.HandlerFunc(pageH.HandleNotFound)

	return r
}
