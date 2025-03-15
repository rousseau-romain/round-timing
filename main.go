package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/helper"
	"github.com/rousseau-romain/round-timing/i18n/locales"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/service/auth"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func languageMiddleware(handler http.Handler, auth *auth.AuthService, slog *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := auth.GetAuthenticateUserFromRequest(r, slog)
		var lang string
		var err error
		if user.Id == 0 {
			lang = helper.GetPreferredLanguage(r)
		} else {
			lang, err = model.GetLanguageLocaleById(user.IdLanguage)
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		ctx, err := ctxi18n.WithLocale(r.Context(), lang)
		if err != nil {
			slog.Error("error setting locale", "error", err)
			http.Error(w, "error setting locale", http.StatusBadRequest)
			return
		}
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
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

	keys := make([]string, 0, len(helper.SupportedLanguages))
	for k := range helper.SupportedLanguages {
		keys = append(keys, k)
	}
	regexCode := strings.Join(keys, "|")

	// PUBLIC ROUTE
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// PAGES ROUTE
	r.Handle("/", auth.AllowToBeAuth(handler.HandlersHome, authService, versionLogger)).Methods("GET")
	r.Handle("/commit-id", auth.AllowToBeAuth(handler.HandlerCommitId, authService, versionLogger)).Methods("GET")
	r.Handle("/version", auth.AllowToBeAuth(handler.HandlerVersion, authService, versionLogger)).Methods("GET")

	r.Handle("/privacy", auth.AllowToBeAuth(handler.HandlerPrivacy, authService, versionLogger)).Methods("GET")
	r.Handle("/cgu", auth.AllowToBeAuth(handler.HandlerCGU, authService, versionLogger)).Methods("GET")
	r.Handle("/match", auth.RequireAuth(handler.HandlersListMatch, authService, versionLogger)).Methods("GET")
	r.Handle("/match", auth.RequireAuth(handler.HandlersCreateMatch, authService, versionLogger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersDeleteMatch, authService, versionLogger)).Methods("DELETE")
	r.Handle("/match/{idMatch:[0-9]+}/unautorized", auth.RequireAuth(handler.HandlersMatchUnAuthorized, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersMatch, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/spectate", auth.RequireAuthAndSpectateOfUserMatch(handler.HandlerSpectateMatch, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/start", auth.RequireAuthAndHisMatch(handler.HandlerStartMatchPage, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/reset", auth.RequireAuthAndHisMatch(handler.HandlerResetMatchPage, authService, versionLogger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/increase-round", auth.RequireAuthAndHisMatch(handler.HandlerMatchNextRound, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/table-live", auth.AllowToBeAuth(handler.HandlerMatchTableLive, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/toggle-mastery/{toggleBool:[0-1]}", auth.RequireAuthAndHisMatch(handler.HandlerToggleMatchMastery, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/use", auth.RequireAuthAndHisMatch(handler.HandlerUsePlayerSpell, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/remove-round-recovery", auth.RequireAuthAndHisMatch(handler.HandlerRemoveRoundRecoveryPlayerSpell, authService, versionLogger)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player", auth.RequireAuthAndHisMatch(handler.HandlersCreatePlayer, authService, versionLogger)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersUpdatePlayer, authService, versionLogger)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersDeletePlayer, authService, versionLogger)).Methods("DELETE")

	r.Handle("/profile", auth.RequireAuth(handler.HandlersProfile, authService, versionLogger)).Methods("GET")
	r.Handle("/profile/configuration/{idConfiguration:[0-9]+}/toggle-configuration", auth.RequireAuth(handler.HandlersProfileToggleUserConfiguration, authService, versionLogger)).Methods("PATCH")
	r.Handle("/profile/spell-favorite/{idSpell:[0-9]+}/toggle-favorite", auth.RequireAuth(handler.HandlersToggleSpellFavorite, authService, versionLogger)).Methods("PATCH")
	r.Handle("/profile/user-spectate", auth.RequireAuth(handler.HandlersProfileAddSpectate, authService, versionLogger)).Methods("POST")
	r.Handle("/profile/user-spectate", auth.RequireAuth(handler.HandlersProfileDeleteSpectate, authService, versionLogger)).Methods("DELETE")
	r.Handle(fmt.Sprintf("/user/{idUser:[0-9]+}/locale/{code:(?:%s)}", regexCode), auth.RequireAuthAndHisAccount(handler.HandlersPlayerLanguage, authService, versionLogger)).Methods("PATCH")

	r.Handle("/signup", auth.RequireNotAuth(handler.HandleSignupEmail, authService, versionLogger)).Methods("GET")
	r.Handle("/signin", auth.RequireNotAuth(handler.HandleLogin, authService, versionLogger)).Methods("GET")

	r.HandleFunc("/signup", handler.HandleCreateUser).Methods("POST")
	r.HandleFunc("/signin", handler.HandleLoginEmail).Methods("POST")
	r.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	r.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")

	r.Handle("/admin/user", auth.RequireAuthAndAdmin(handler.HandlersListUser, authService, versionLogger)).Methods("GET")
	r.Handle("/admin/user/{idUser:[0-9]+}/toggle-enabled/{toggleEnabled:(?:true|false)}", auth.RequireAuthAndAdmin(handler.HandlersUserEnabled, authService, versionLogger)).Methods("PATCH")

	r.Handle("/404", auth.AllowToBeAuth(handler.HandlersNotFound, authService, versionLogger)).Methods("GET")
	r.Handle("/403", auth.AllowToBeAuth(handler.HandlersForbidden, authService, versionLogger)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(handler.HandlersNotFound)

	return http.ListenAndServe(":2468", languageMiddleware(r, authService, versionLogger))
}
