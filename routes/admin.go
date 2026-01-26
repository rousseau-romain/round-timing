package routes

import (
	"log/slog"

	"github.com/gorilla/mux"
	handlersAdmin "github.com/rousseau-romain/round-timing/handlers/admin"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerAdminRoutes(r *mux.Router, handler *handlersAdmin.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/admin/user", middleware.RequireAuthAndAdmin(handler.HandleListUser, authService, logger)).Methods("GET")
	r.Handle("/admin/user/{idUser:[0-9]+}/toggle-enabled/{toggleEnabled:(?:true|false)}", middleware.RequireAuthAndAdmin(handler.HandleUserEnabled, authService, logger)).Methods("PATCH")
}
