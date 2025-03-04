package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/shared/components"
)

type Handler struct {
	auth      *auth.AuthService
	error     components.PopinMessages
	languages []model.Language
}

func New(auth *auth.AuthService) *Handler {
	languages, err := model.GetLanguages()
	if err != nil {
		log.Println(err)
		languages = []model.Language{}
	}
	return &Handler{
		auth: auth,
		error: components.PopinMessages{
			Title:    "",
			Messages: []string{},
		},
		languages: languages,
	}
}

func RenderComponentErrorAndLog(title string, message, messageLog []string, httpCode int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(httpCode)
	components.PopinMessage(components.PopinMessages{
		Title:    title,
		Messages: message,
		Type:     "error",
	}).Render(r.Context(), w)
	log.Println(strings.Join(messageLog, "\n"))
}
