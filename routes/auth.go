package routes

import (
	"log/slog"
	"time"

	"github.com/gorilla/mux"
	handlersAuth "github.com/rousseau-romain/round-timing/handlers/auth"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/service/auth"
)

var authLimiter = middleware.NewRateLimiter(10, 15*time.Minute)

func registerAuthRoutes(r *mux.Router, handler *handlersAuth.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/signup", middleware.RequireNotAuth(handler.HandleSignupEmail, authService, logger)).Methods("GET")
	r.Handle("/signin", middleware.RequireNotAuth(handler.HandleLogin, authService, logger)).Methods("GET")

	r.HandleFunc("/signup", authLimiter.Limit(handler.HandleCreateUser)).Methods("POST")
	r.HandleFunc("/signin", authLimiter.Limit(handler.HandleLoginEmail)).Methods("POST")
	r.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	r.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")
}
