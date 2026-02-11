package config

import (
	"log"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/joho/godotenv"
)

// Set via -ldflags at build time, used as fallback when VCS info is unavailable
var (
	version     = ""
	commit      = ""
	buildTime   = ""
	vcsModified = ""
)

var (
	VERSION,
	COMMIT,
	BUILD_TIME,
	VCS_MODIFIED,
	DB_HOST,
	DB_NAME,
	DB_PASSWORD,
	DB_USER,
	DB_URL,
	DB_DRIVER,
	SALT_SECRET,
	JWT_SECRET_KEY,
	DB_PORT,
	PUBLIC_HOST_PORT,
	COOKIES_AUTH_SECRET,
	CSRF_KEY,
	DISCORD_CLIENT_ID,
	DISCORD_CLIENT_SECRET,
	GOOGLE_CLIENT_ID,
	GOOGLE_CLIENT_SECRET,
	GITHUB_CLIENT_ID,
	GITHUB_CLIENT_SECRET,
	ENV string
	COOKIES_AUTH_AGE_IN_SECONDS int
	COOKIES_AUTH_IS_SECURE      bool
	COOKIES_AUTH_IS_HTTP_ONLY   bool
)

func requiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}

func init() {
	log.SetFlags(log.Llongfile)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	VERSION = "v1.7.0"
	COMMIT = "unknown"
	BUILD_TIME = "unknown"
	VCS_MODIFIED = "false"

	// Try debug.ReadBuildInfo() first (works with go build locally)
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				COMMIT = s.Value
			case "vcs.time":
				BUILD_TIME = s.Value
			case "vcs.modified":
				VCS_MODIFIED = s.Value
			}
		}
	}

	// Fall back to ldflags values (set during Docker build)
	if COMMIT == "unknown" && commit != "" {
		COMMIT = commit
	}
	if BUILD_TIME == "unknown" && buildTime != "" {
		BUILD_TIME = buildTime
	}
	if VCS_MODIFIED == "false" && vcsModified != "" {
		VCS_MODIFIED = vcsModified
	}
	if version != "" {
		VERSION = version
	}

	ENV = os.Getenv("ENV")

	DB_HOST = requiredEnv("DB_HOST")
	DB_NAME = requiredEnv("DB_NAME")
	DB_PASSWORD = requiredEnv("DB_PASSWORD")
	DB_USER = requiredEnv("DB_USER")
	DB_DRIVER = requiredEnv("DB_DRIVER")
	DB_PORT = requiredEnv("DB_PORT")
	DB_URL = os.Getenv("DB_URL")

	SALT_SECRET = requiredEnv("SALT_SECRET")
	JWT_SECRET_KEY = requiredEnv("JWT_SECRET_KEY")
	PUBLIC_HOST_PORT = requiredEnv("PUBLIC_HOST_PORT")
	COOKIES_AUTH_SECRET = requiredEnv("COOKIES_AUTH_SECRET")
	CSRF_KEY = requiredEnv("CSRF_KEY")

	COOKIES_AUTH_AGE_IN_SECONDS = 60 * 60 * 24 * 2
	COOKIES_AUTH_IS_SECURE, _ = strconv.ParseBool(os.Getenv("COOKIES_AUTH_IS_SECURE"))
	COOKIES_AUTH_IS_HTTP_ONLY, _ = strconv.ParseBool(os.Getenv("COOKIES_AUTH_IS_HTTP_ONLY"))
	DISCORD_CLIENT_ID = os.Getenv("DISCORD_CLIENT_ID")
	DISCORD_CLIENT_SECRET = os.Getenv("DISCORD_CLIENT_SECRET")
	GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
	GITHUB_CLIENT_ID = os.Getenv("GITHUB_CLIENT_ID")
	GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
}
