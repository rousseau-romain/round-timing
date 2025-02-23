package main

import (
	"fmt"
	"log"
	"net/http"
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

func languageMiddleware(handler http.Handler, auth *auth.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		var lang string
		if err != nil {
			lang = helper.GetPreferredLanguage(r)
		} else {
			user, err := model.GetUserByOauth2Id(session.UserID)
			if user.Id != 0 {
				lang, err = model.GetLanguageLocaleById(user.IdLanguage)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		ctx, err := ctxi18n.WithLocale(r.Context(), lang)
		if err != nil {
			log.Printf("error setting locale: %v", err)
			http.Error(w, "error setting locale", http.StatusBadRequest)
			return
		}
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func run() error {
	if err := ctxi18n.Load(locales.Content); err != nil {
		log.Printf("error loading locales: %v", err)
	}

	sessionStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: config.COOKIES_AUTH_SECRET,
		MaxAge:     config.COOKIES_AUTH_AGE_IN_SECONDS,
		Secure:     config.COOKIES_AUTH_IS_SECURE,
		HttpOnly:   config.COOKIES_AUTH_IS_HTTP_ONLY,
	})
	authService := auth.NewAuthService(sessionStore)

	r := mux.NewRouter()

	handler := handlers.New(authService)

	keys := make([]string, 0, len(helper.SupportedLanguages))
	for k := range helper.SupportedLanguages {
		keys = append(keys, k)
	}
	regexCode := strings.Join(keys, "|")

	// PUBLIC ROUTE
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// PAGES ROUTE
	r.Handle("/", auth.AllowToBeAuth(handler.HandlersHome, authService)).Methods("GET")
	r.Handle("/commit-id", auth.AllowToBeAuth(handler.HandlerCommitId, authService)).Methods("GET")

	r.Handle("/privacy", auth.AllowToBeAuth(handler.HandlerPrivacy, authService)).Methods("GET")
	r.Handle("/cgu", auth.AllowToBeAuth(handler.HandlerCGU, authService)).Methods("GET")
	r.Handle("/match", auth.RequireAuth(handler.HandlersListMatch, authService)).Methods("GET")
	r.Handle("/match", auth.RequireAuth(handler.HandlersCreateMatch, authService)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersDeleteMatch, authService)).Methods("DELETE")
	r.Handle("/match/{idMatch:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersMatch, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/spectate", auth.RequireAuthAndSpectateOfUserMatch(handler.HandlerSpectateMatch, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/start", auth.RequireAuthAndHisMatch(handler.HandlerStartMatchPage, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/reset", auth.RequireAuthAndHisMatch(handler.HandlerResetMatchPage, authService)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/increase-round", auth.RequireAuthAndHisMatch(handler.HandlerMatchNextRound, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/table-live", auth.AllowToBeAuth(handler.HandlerMatchTableLive, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/toggle-mastery/{toggleBool:[0-1]}", auth.RequireAuthAndHisMatch(handler.HandlerToggleMatchMastery, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/use", auth.RequireAuthAndHisMatch(handler.HandlerUsePlayerSpell, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/remove-round-recovery", auth.RequireAuthAndHisMatch(handler.HandlerRemoveRoundRecoveryPlayerSpell, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player", auth.RequireAuthAndHisMatch(handler.HandlersCreatePlayer, authService)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersUpdatePlayer, authService)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersDeletePlayer, authService)).Methods("DELETE")

	r.Handle("/profile", auth.RequireAuth(handler.HandlersProfile, authService)).Methods("GET")
	r.Handle("/profile/spell-favorite/{idSpell:[0-9]+}/toggle-favorite", auth.RequireAuth(handler.HandlersToggleSpellFavorite, authService)).Methods("PATCH")
	r.Handle("/profile/user-spectate", auth.RequireAuth(handler.HandlersProfileAddSpectate, authService)).Methods("POST")
	r.Handle("/profile/user-spectate", auth.RequireAuth(handler.HandlersProfileDeleteSpectate, authService)).Methods("DELETE")
	r.Handle(fmt.Sprintf("/user/{idUser:[0-9]+}/locale/{code:(?:%s)}", regexCode), auth.RequireAuthAndHisAccount(handler.HandlersPlayerLanguage, authService)).Methods("PATCH")

	r.Handle("/signin", auth.RequireNotAuth(handler.HandleLogin, authService)).Methods("GET")
	r.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	r.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")

	r.Handle("/admin/user", auth.RequireAuthAndAdmin(handler.HandlersListUser, authService)).Methods("GET")
	r.Handle("/admin/user/{idUser:[0-9]+}/toggle-enabled/{toggleEnabled:(?:true|false)}", auth.RequireAuthAndAdmin(handler.HandlersUserEnabled, authService)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(handler.HandlersNotFound)

	return http.ListenAndServe(":2468", languageMiddleware(r, authService))
}
