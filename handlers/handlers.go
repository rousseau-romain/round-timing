package handlers

import (
	"log"
	"net/http"

	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/shared/components"
)

type Handler struct {
	auth      *auth.AuthService
	error     components.Error
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
		error: components.Error{
			Title:    "",
			Messages: []string{},
		},
		languages: languages,
	}
}

func RenderComponentError(title string, message []string, w http.ResponseWriter, r *http.Request) {
	components.ErrorMessages(components.Error{
		Title:    title,
		Messages: message,
	}).Render(r.Context(), w)
}
