package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/invopop/ctxi18n/i18n"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	"github.com/rousseau-romain/round-timing/model/system"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/views/components/layout"
)

type Handler struct {
	Auth      *auth.AuthService
	Slog      *slog.Logger
	Error     layout.PopinMessages
	Languages []system.Language
}

func New(auth *auth.AuthService, slog *slog.Logger) *Handler {
	languages, err := system.GetLanguages(context.Background())
	if err != nil {
		languages = []system.Language{}
	}
	return &Handler{
		Auth: auth,
		Slog: slog,
		Error: layout.PopinMessages{
			Title:    "",
			Messages: []string{},
		},
		Languages: languages,
	}
}

func GetPageNavDefault(r *http.Request) []layout.NavItem {
	return []layout.NavItem{
		{
			Name: i18n.T(r.Context(), "page.match-list.title"),
			Url:  "match",
		},
	}
}

func (h *Handler) GetPageNavCustom(r *http.Request, user userModel.User, match matchModel.Match) []layout.NavItem {
	var pageNav = GetPageNavDefault(r)
	if user.Id != 0 {
		if match.Id != 0 {
			pageNav = append(pageNav, layout.NavItem{
				Name: fmt.Sprintf("Match %s (%d)", match.Name, match.Id),
				Url:  fmt.Sprintf("match/%d", match.Id),
			})
		} else {
			lastMatch, err := matchModel.GetLastMatchByUserId(r.Context(), user.Id)
			if err != nil {
				h.Slog.Error(err.Error())
				return pageNav
			}
			if lastMatch.Id != 0 {
				pageNav = append(pageNav, layout.NavItem{
					Name: i18n.T(r.Context(), "global.header.last-match", i18n.M{"name": lastMatch.Name, "id": lastMatch.Id}),
					Url:  fmt.Sprintf("match/%d", lastMatch.Id),
				})
			}
		}
	}
	if user.IsAdmin {
		pageNav = append(pageNav, layout.NavItem{
			Name: "Admin list users",
			Url:  "admin/user",
		})
	}
	return pageNav
}

func RespondWithError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error, userMessage string, status int) {
	logger.Error(err.Error())
	if r.Header.Get("HX-Request") == "true" {
		RenderComponentError(userMessage, []string{userMessage}, status, w, r)
	} else {
		http.Error(w, userMessage, status)
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
