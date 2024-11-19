package main

import (
	"log"
	"net/http"

	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/service/auth"

	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	sessionStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: config.COOKIES_AUTH_SECRET,
		MaxAge:     config.COOKIES_AUTH_AGE_IN_SECONDS,
		Secure:     config.COOKIES_AUTH_IS_SECURE,
		HttpOnly:   config.COOKIES_AUTH_IS_HTTP_ONLY,
	})
	authService := auth.NewAuthService(sessionStore)

	r := mux.NewRouter()

	handler := handlers.New(authService)

	// PUBLIC ROUTE
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// PAGES ROUTE
	r.Handle("/", auth.AllowToBeAuth(handler.HandlersHome, authService)).Methods("GET")

	r.Handle("/match", auth.RequireAuth(handler.HandlersListMatch, authService)).Methods("GET")
	r.Handle("/match", auth.RequireAuth(handler.HandlersCreateMatch, authService)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersDeleteMatch, authService)).Methods("DELETE")
	r.Handle("/match/{idMatch:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersMatch, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/start", auth.RequireAuthAndHisMatch(handler.HandlerStartMatchPage, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/reset", auth.RequireAuthAndHisMatch(handler.HandlerResetMatchPage, authService)).Methods("PATCH")
	r.Handle("/match/{idMatch:[0-9]+}/increase-round", auth.RequireAuthAndHisMatch(handler.HandlerMatchNextRound, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/toggle-mastery/{toggleBool:[0-1]}", auth.RequireAuthAndHisMatch(handler.HandlerToggleMatchMastery, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/use", auth.RequireAuthAndHisMatch(handler.HandlerUsePlayerSpell, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player-spell/{idPlayerSpell:[0-9]+}/remove-round-recovery", auth.RequireAuthAndHisMatch(handler.HandlerRemoveRoundRecoveryPlayerSpell, authService)).Methods("GET")
	r.Handle("/match/{idMatch:[0-9]+}/player", auth.RequireAuthAndHisMatch(handler.HandlersCreatePlayer, authService)).Methods("POST")
	r.Handle("/match/{idMatch:[0-9]+}/player/{idPlayer:[0-9]+}", auth.RequireAuthAndHisMatch(handler.HandlersDeletePlayer, authService)).Methods("DELETE")

	r.Handle("/profile", auth.RequireAuth(handler.HandlersProfile, authService)).Methods("GET")

	r.Handle("/signin", auth.RequireNotAuth(handler.HandleLogin, authService)).Methods("GET")
	r.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	r.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")

	r.Handle("/admin/user", auth.RequireAuthAndAdmin(handler.HandlersListUser, authService)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(handler.HandlersNotFound)
	return http.ListenAndServe(":2468", r)
}
