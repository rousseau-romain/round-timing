package helper

import (
	"fmt"
	"net/http"
	"strings"
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
