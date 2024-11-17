package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DB_HOST,
	DB_NAME,
	DB_PASSWORD,
	DB_USER,
	DB_URL,
	DB_DRIVER,
	JWT_SECRET_KEY,
	DB_PORT,
	PUBLIC_HOST_PORT,
	COOKIES_AUTH_SECRET,
	DISCORD_CLIENT_ID,
	DISCORD_CLIENT_SECRET,
	GITHUB_CLIENT_ID,
	GITHUB_CLIENT_SECRET string
	COOKIES_AUTH_AGE_IN_SECONDS int
	COOKIES_AUTH_IS_SECURE      bool
	COOKIES_AUTH_IS_HTTP_ONLY   bool
)

func init() {
	log.SetFlags(log.Llongfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Coudn't load env file!!")
	}

	DB_HOST = os.Getenv("DB_HOST")
	DB_NAME = os.Getenv("DB_NAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_USER = os.Getenv("DB_USER")
	DB_DRIVER = os.Getenv("DB_DRIVER")
	DB_PORT = os.Getenv("DB_PORT")
	DB_URL = os.Getenv("DB_URL")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

	PUBLIC_HOST_PORT = os.Getenv("PUBLIC_HOST_PORT")
	COOKIES_AUTH_SECRET = os.Getenv("COOKIES_AUTH_SECRET")

	COOKIES_AUTH_AGE_IN_SECONDS = 60 * 60 * 24 * 2
	COOKIES_AUTH_IS_SECURE, _ = strconv.ParseBool(os.Getenv("COOKIES_AUTH_IS_SECURE"))
	COOKIES_AUTH_IS_HTTP_ONLY, _ = strconv.ParseBool(os.Getenv("COOKIES_AUTH_IS_HTTP_ONLY"))
	DISCORD_CLIENT_ID = os.Getenv("DISCORD_CLIENT_ID")
	DISCORD_CLIENT_SECRET = os.Getenv("DISCORD_CLIENT_SECRET")
	GITHUB_CLIENT_ID = os.Getenv("GITHUB_CLIENT_ID")
	GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
}
