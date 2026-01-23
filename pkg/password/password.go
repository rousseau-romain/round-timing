package password

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"unicode"

	"github.com/invopop/ctxi18n/i18n"
	"golang.org/x/crypto/argon2"
)

// GenerateSalt generates a random 16-byte salt encoded as base64.
func GenerateSalt() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Hash hashes a password with the given salt using Argon2id.
func Hash(password, salt string) string {
	hash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash) + ":" + salt
}

// Check verifies a password against a stored hash.
func Check(storedHash, password string) bool {
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return false
	}
	hash := Hash(password, parts[1])
	return storedHash == hash
}

// Validate checks if a password meets security requirements.
// Returns true if valid, or false with a list of error messages.
func Validate(r *http.Request, password string) (bool, []string) {
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
		errors = append(errors, i18n.T(r.Context(), "page.signup.error.password.too-short"))
	}
	if !hasUpper {
		errors = append(errors, i18n.T(r.Context(), "page.signup.error.password.uppercase"))
	}
	if !hasLower {
		errors = append(errors, i18n.T(r.Context(), "page.signup.error.password.lowercase"))
	}
	if !hasNumber {
		errors = append(errors, i18n.T(r.Context(), "page.signup.error.password.digit"))
	}
	if !hasSpecial {
		errors = append(errors, i18n.T(r.Context(), "page.signup.error.password.special"))
	}

	if len(errors) == 0 {
		return true, nil
	}

	return false, errors
}
