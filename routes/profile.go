package routes

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gorilla/mux"
	handlersProfile "github.com/rousseau-romain/round-timing/handlers/profile"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/pkg/lang"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func registerProfileRoutes(r *mux.Router, handler *handlersProfile.Handler, authService *auth.AuthService, logger *slog.Logger) {
	r.Handle("/profile", middleware.RequireAuth(handler.HandleProfile, authService, logger)).Methods("GET")
	r.Handle("/profile/configuration/{idConfiguration:[0-9]+}/toggle-configuration", middleware.RequireAuth(handler.HandleProfileToggleUserConfiguration, authService, logger)).Methods("PATCH")
	r.Handle("/profile/spell-favorite/{idSpell:[0-9]+}/toggle-favorite", middleware.RequireAuth(handler.HandleToggleSpellFavorite, authService, logger)).Methods("PATCH")
	r.Handle("/profile/user-spectate", middleware.RequireAuth(handler.HandleProfileAddSpectate, authService, logger)).Methods("POST")
	r.Handle("/profile/user-spectate", middleware.RequireAuth(handler.HandleProfileDeleteSpectate, authService, logger)).Methods("DELETE")

	// User locale route
	keys := make([]string, 0, len(lang.SupportedLanguages))
	for k := range lang.SupportedLanguages {
		keys = append(keys, k)
	}
	regexCode := strings.Join(keys, "|")
	r.Handle(fmt.Sprintf("/user/{idUser:[0-9]+}/locale/{code:(?:%s)}", regexCode), middleware.RequireAuthAndHisAccount(handler.HandlePlayerLanguage, authService, logger)).Methods("PATCH")
}
