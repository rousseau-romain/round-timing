package handlers

import (
	"log"

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
