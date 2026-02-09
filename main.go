package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/i18n/locales"
	"github.com/rousseau-romain/round-timing/middleware"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/routes"
	"github.com/rousseau-romain/round-timing/service/auth"

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

	// Log DB pool stats every 5 minutes
	model.StartDBStatsLogger(5 * time.Minute)

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

	handler := handlers.New(authService, versionLogger)

	router := routes.Setup(handler, authService, versionLogger)

	return http.ListenAndServe(":2468", middleware.Language(router, authService, versionLogger))
}
