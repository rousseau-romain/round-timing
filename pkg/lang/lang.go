package lang

import (
	"net/http"
	"strings"
)

// SupportedLanguages maps locale codes to their database IDs.
var SupportedLanguages = map[string]int{
	"en": 1,
	"fr": 2,
	"it": 3,
	"es": 4,
	"pt": 5,
}

// GetPreferred detects the preferred language from the Accept-Language header.
// Falls back to "en" if no supported language is found.
func GetPreferred(r *http.Request) string {
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
