package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
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

	csrfMiddleware := csrf.Protect(
		[]byte(config.CSRF_KEY),
		csrf.Secure(config.COOKIES_AUTH_IS_SECURE),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			versionLogger.Warn("CSRF validation failed",
				"method", r.Method,
				"path", r.URL.Path,
				"reason", csrf.FailureReason(r),
				"hasCSRFHeader", r.Header.Get("X-CSRF-Token") != "",
			)
			http.Error(w, "Forbidden - invalid CSRF token", http.StatusForbidden)
		})),
	)

	inner := csrfMiddleware(middleware.WithCSRFToken(middleware.Language(router, authService, versionLogger)))

	// gorilla/csrf v1.7.3 defaults to assuming HTTPS. When running on plain
	// HTTP (Secure=false), we must signal this via PlaintextHTTPRequest so
	// the Origin/Referer checks use the correct scheme.
	var chain http.Handler = inner
	if !config.COOKIES_AUTH_IS_SECURE {
		chain = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inner.ServeHTTP(w, csrf.PlaintextHTTPRequest(r))
		})
	}

	return http.ListenAndServe(":2468", chain)
}
