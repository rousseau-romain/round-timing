package handlers

import (
	"log/slog"
	"net/http"

	"github.com/rousseau-romain/round-timing/model/system"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/views/components/layout"
)

type Handler struct {
	auth      *auth.AuthService
	Slog      *slog.Logger
	error     layout.PopinMessages
	languages []system.Language
}

func New(auth *auth.AuthService, slog *slog.Logger) *Handler {
	languages, err := system.GetLanguages()
	if err != nil {
		languages = []system.Language{}
	}
	return &Handler{
		auth: auth,
		Slog: slog,
		error: layout.PopinMessages{
			Title:    "",
			Messages: []string{},
		},
		languages: languages,
	}
}

func RenderComponentError(title string, message []string, httpCode int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(httpCode)
	layout.PopinMessage(layout.PopinMessages{
		Title:    title,
		Messages: message,
		Type:     "error",
	}).Render(r.Context(), w)
}

func RenderComponentWarning(title string, message []string, httpCode int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(httpCode)
	layout.PopinMessage(layout.PopinMessages{
		Title:    title,
		Messages: message,
		Type:     "warning",
	}).Render(r.Context(), w)
}
