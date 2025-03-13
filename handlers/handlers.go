package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/shared/components"
)

type Handler struct {
	auth      *auth.AuthService
	Slog      *slog.Logger
	error     components.PopinMessages
	languages []model.Language
}

func New(auth *auth.AuthService, slog *slog.Logger) *Handler {
	languages, err := model.GetLanguages()
	if err != nil {
		languages = []model.Language{}
	}
	return &Handler{
		auth: auth,
		Slog: slog,
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
	slog.Error(strings.Join(messageLog, "\n"))
}
