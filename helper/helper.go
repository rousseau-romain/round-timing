package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/invopop/ctxi18n/i18n"
	"golang.org/x/crypto/argon2"
)

var MailContact = "roundtiming@gmail.com"

// Id in database
var SupportedLanguages = map[string]int{
	"en": 1,
	"fr": 2,
	"it": 3,
	"es": 4,
	"pt": 5,
}

func GetPreferredLanguage(r *http.Request) string {
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang == "" {
		return "en"
	}

	languages := strings.Split(acceptLang, ",")
	for _, lang := range languages {
		langCode := strings.TrimSpace(strings.Split(lang, ";")[0])
		langCode = strings.Split(langCode, "-")[0]

		if _, ok := SupportedLanguages[langCode]; ok {
			return langCode
		}
	}

	return "en"
}

var logTypeTemplate = func(e any) string {
	fmt.Printf("Type %T\n", e)
	return ""
}

// 139 is mastery hammer
var MasteryIdSpells = []int{136, 137, 138, 140, 141, 142, 143}

func GenerateSalt() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HashPassword(password, salt string) string {
	hash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash) + ":" + salt
}

func CheckPassword(storedHash, password string) bool {
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return false
	}
	hash := HashPassword(password, parts[1])
	return storedHash == hash
}

func IsValidPassword(r *http.Request, password string) (bool, []string) {
	var errors []string
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		errors = append(errors, i18n.T(r.Context(), "page.signup.password.error.too-short"))
	}
	if !hasUpper {
		errors = append(errors, i18n.T(r.Context(), "page.signup.password.error.uppercase"))
	}
	if !hasLower {
		errors = append(errors, i18n.T(r.Context(), "page.signup.password.error.lowercase"))
	}
	if !hasNumber {
		errors = append(errors, i18n.T(r.Context(), "page.signup.password.error.digit"))
	}
	if !hasSpecial {
		errors = append(errors, i18n.T(r.Context(), "page.signup.password.error.special"))
	}

	if len(errors) == 0 {
		return true, nil
	}

	return false, errors
}
