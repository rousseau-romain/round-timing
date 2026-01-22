package handlers

import (
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"strings"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/components/layout"
	"github.com/rousseau-romain/round-timing/views/page"
	"github.com/rousseau-romain/round-timing/views/page/legal"

	"io"
)

func GetPageNavDefault(r *http.Request) []layout.NavItem {
	return []layout.NavItem{
		{
			Name: i18n.T(r.Context(), "page.match-list.title"),
			Url:  "match",
		},
	}
}

func (h *Handler) GetPageNavCustom(r *http.Request, user model.User, match model.Match) []layout.NavItem {
	var pageNav = GetPageNavDefault(r)
	if user.Id != 0 {
		if match.Id != 0 {
			pageNav = append(pageNav, layout.NavItem{
				Name: fmt.Sprintf("Match %s (%d)", match.Name, match.Id),
				Url:  fmt.Sprintf("match/%d", match.Id),
			})
		} else {
			lastMatch, err := model.GetLastMatchByUserId(user.Id)
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
func (h *Handler) HandlerCommitId(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("config/commit-id.json")
	if err != nil {
		h.Slog.Error(err.Error())
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var commit struct {
		CommitId string `json:"commit_id"`
	}

	json.Unmarshal(byteValue, &commit)

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

func (h *Handler) HandlerVersion(w http.ResponseWriter, r *http.Request) {
	var version struct {
		Version string `json:"version"`
	}

	version.Version = config.VERSION

	byteValue, _ := json.Marshal(version)

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

func (h *Handler) HandlersNotFound(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	page.NotFoundPage("", h.GetPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path, user).Render(r.Context(), w)
}

func (h *Handler) HandlersForbidden(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	page.ForbidenPage("", h.GetPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path, user).Render(r.Context(), w)
}

func (h *Handler) HandlersHome(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	pageNav := GetPageNavDefault(r)

	h.error = layout.PopinMessages{
		Title:    r.URL.Query().Get("errorTitle"),
		Messages: strings.Split(r.URL.Query().Get("errorMessages"), ","),
	}

	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
		pageNav = h.GetPageNavCustom(r, user, model.Match{})
		page.HomePage(user, h.error, pageNav, h.languages, r.URL.Path).Render(r.Context(), w)
		return
	}
	page.HomePage(model.User{}, h.error, pageNav, h.languages, r.URL.Path).Render(r.Context(), w)
}

func (h *Handler) HandlerCGU(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	legal.CGU(h.error, h.GetPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path).Render(r.Context(), w)
}

func (h *Handler) HandlerPrivacy(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	legal.Privacy(h.error, h.GetPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path).Render(r.Context(), w)
}
