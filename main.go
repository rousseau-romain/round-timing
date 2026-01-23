package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/i18n/locales"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/pkg/lang"
	"github.com/rousseau-romain/round-timing/service/auth"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n"
)

func init() {
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic("failed to load Europe/Paris timezone: " + err.Error())
	}
	time.Local = loc
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)
	versionLogger := logger.With("version", config.VERSION)

	if err := ctxi18n.Load(locales.Content); err != nil {
		versionLogger.Error("error loading locales", "error", err)
	}

	sessionStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: config.COOKIES_AUTH_SECRET,
		MaxAge:     config.COOKIES_AUTH_AGE_IN_SECONDS,
		Secure:     config.COOKIES_AUTH_IS_SECURE,
		HttpOnly:   config.COOKIES_AUTH_IS_HTTP_ONLY,
	})
	authService := auth.NewAuthService(sessionStore)

	r := mux.NewRouter()

	handler := handlers.New(authService, versionLogger)

	keys := make([]string, 0, len(lang.SupportedLanguages))
	for k := range lang.SupportedLanguages {
		keys = append(keys, k)
	}
	regexCode := strings.Join(keys, "|")

	// PUBLIC ROUTE
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// PAGES ROUTE
	r.Handle("/", middleware.AllowToBeAuth(handler.HandlersHome, authService, versionLogger)).Methods("GET")
	r.Handle("/commit-id", middleware.AllowToBeAuth(handler.HandlerCommitId, authService, versionLogger)).Methods("GET")
	r.Handle("/version", middleware.AllowToBeAuth(handler.HandlerVersion, authService, versionLogger)).Methods("GET")

	r.Handle("/privacy", middleware.AllowToBeAuth(handler.HandlerPrivacy, authService, versionLogger)).Methods("GET")
	r.Handle("/cgu", middleware.AllowToBeAuth(handler.HandlerCGU, authService, versionLogger)).Methods("GET")
	r.Handle("/match", middleware.RequireAuth(handler.HandlersListMatch, authService, versionLogger)).Methods("GET")
	r.Handle("/match", middleware.RequireAuth(handler.HandlersCreateMatch, authService, versionLogger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersDeleteMatch, authService, versionLogger)).Methods("DELETE")
	r.Handle("/match/{idMatch:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersMatch, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/spectate", middleware.RequireAuthAndSpectateOfUserMatch(handler.HandlerSpectateMatch, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/start", middleware.RequireAuthAndHisMatch(handler.HandlerStartMatchPage, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/reset", middleware.RequireAuthAndHisMatch(handler.HandlerResetMatchPage, authService, versionLogger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/increase-round", middleware.RequireAuthAndHisMatch(handler.HandlerMatchNextRound, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/table-live", middleware.AllowToBeAuth(handler.HandlerMatchTableLive, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/toggle-mastery/{toggleBool:[0-1]}", middleware.RequireAuthAndHisMatch(handler.HandlerToggleMatchMastery, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/use", middleware.RequireAuthAndHisMatch(handler.HandlerUsePlayerSpell, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/remove-round-recovery", middleware.RequireAuthAndHisMatch(handler.HandlerRemoveRoundRecoveryPlayerSpell, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player", middleware.RequireAuthAndHisMatch(handler.HandlersCreatePlayer, authService, versionLogger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersUpdatePlayer, authService, versionLogger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", middleware.RequireAuthAndHisMatch(handler.HandlersDeletePlayer, authService, versionLogger)).Methods("DELETE")

	r.Handle("/profile", middleware.RequireAuth(handler.HandlersProfile, authService, versionLogger)).Methods("GET")
	r.Handle("/profile/configuration/{idConfiguration:[0-9]+}/toggle-configuration", middleware.RequireAuth(handler.HandlersProfileToggleUserConfiguration, authService, versionLogger)).Methods("PATCH")
	r.Handle("/profile/spell-favorite/{idSpell:[0-9]+}/toggle-favorite", middleware.RequireAuth(handler.HandlersToggleSpellFavorite, authService, versionLogger)).Methods("PATCH")
	r.Handle("/profile/user-spectate", middleware.RequireAuth(handler.HandlersProfileAddSpectate, authService, versionLogger)).Methods("POST")
	r.Handle("/profile/user-spectate", middleware.RequireAuth(handler.HandlersProfileDeleteSpectate, authService, versionLogger)).Methods("DELETE")
	r.Handle(fmt.Sprintf("/user/{idUser:[0-9]+}/locale/{code:(?:%s)}", regexCode), middleware.RequireAuthAndHisAccount(handler.HandlersPlayerLanguage, authService, versionLogger)).Methods("PATCH")

	r.Handle("/signup", middleware.RequireNotAuth(handler.HandleSignupEmail, authService, versionLogger)).Methods("GET")
	r.Handle("/signin", middleware.RequireNotAuth(handler.HandleLogin, authService, versionLogger)).Methods("GET")

	r.HandleFunc("/signup", handler.HandleCreateUser).Methods("POST")
	r.HandleFunc("/signin", handler.HandleLoginEmail).Methods("POST")
	r.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	r.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")

	r.Handle("/admin/user", middleware.RequireAuthAndAdmin(handler.HandlersListUser, authService, versionLogger)).Methods("GET")
	r.Handle("/admin/user/{idUser:[0-9]+}/toggle-enabled/{toggleEnabled:(?:true|false)}", middleware.RequireAuthAndAdmin(handler.HandlersUserEnabled, authService, versionLogger)).Methods("PATCH")

	r.Handle("/404", middleware.AllowToBeAuth(handler.HandlersNotFound, authService, versionLogger)).Methods("GET")
	r.Handle("/403", middleware.AllowToBeAuth(handler.HandlersForbidden, authService, versionLogger)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(handler.HandlersNotFound)

	return http.ListenAndServe(":2468", middleware.Language(r, authService, versionLogger))
}
