package handlers

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/shared/components"
	"github.com/rousseau-romain/round-timing/views/page"

	"io"
)

func GetPageNavDefault(r *http.Request) []components.NavItem {
	return []components.NavItem{
		{
			Name: i18n.T(r.Context(), "page.match-list.title"),
			Url:  "match",
		},
	}
}

func getPageNavCustom(r *http.Request, user model.User, match model.Match) []components.NavItem {
	var pageNav = GetPageNavDefault(r)
	if user.Id != 0 {
		if match.Id != 0 {
			pageNav = append(pageNav, components.NavItem{
				Name: fmt.Sprintf("Match %s (%d)", match.Name, match.Id),
				Url:  fmt.Sprintf("match/%d", match.Id),
			})
		} else {
			lastMatch, err := model.GetLastMatchByUserId(user.Id)
			if err != nil {
				log.Println(err)
				return pageNav
			}
			if lastMatch.Id != 0 {
				pageNav = append(pageNav, components.NavItem{
					Name: i18n.T(r.Context(), "global.header.last-match", i18n.M{"name": lastMatch.Name, "id": lastMatch.Id}),
					Url:  fmt.Sprintf("match/%d", lastMatch.Id),
				})
			}
		}
	}
	if user.IsAdmin {
		pageNav = append(pageNav, components.NavItem{
			Name: "Admin list users",
			Url:  "admin/user",
		})
	}
	return pageNav
}
func (h *Handler) HandlerCommitId(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("config/commit-id.json")
	if err != nil {
		log.Println(err)
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

func (h *Handler) HandlersNotFound(w http.ResponseWriter, r *http.Request) {
	page.NotFoundPage(h.error, GetPageNavDefault(r), h.languages, r.URL.Path).Render(r.Context(), w)
}

func (h *Handler) HandlersHome(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r)
	pageNav := GetPageNavDefault(r)

	h.error = components.PopinMessages{
		Title:    r.URL.Query().Get("errorTitle"),
		Messages: strings.Split(r.URL.Query().Get("errorMessages"), ","),
	}

	if user.Id != 0 {
		pageNav = getPageNavCustom(r, user, model.Match{})
		page.HomePage(user, h.error, pageNav, h.languages, r.URL.Path).Render(r.Context(), w)
		return
	}
	page.HomePage(model.User{}, h.error, pageNav, h.languages, r.URL.Path).Render(r.Context(), w)
}

func (h *Handler) HandlerCGU(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r)
	page.CGU(h.error, getPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path).Render(r.Context(), w)
}

func (h *Handler) HandlerPrivacy(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r)
	page.Privacy(h.error, getPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path).Render(r.Context(), w)
}
