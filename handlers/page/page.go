package page

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/handlers"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/views/components/layout"
	pageView "github.com/rousseau-romain/round-timing/views/page"
	"github.com/rousseau-romain/round-timing/views/page/legal"
)

type Handler struct {
	*handlers.Handler
}

func (h *Handler) HandleCommitId(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) HandleVersion(w http.ResponseWriter, r *http.Request) {
	var version struct {
		Version string `json:"version"`
	}

	version.Version = config.VERSION

	byteValue, _ := json.Marshal(version)

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
}

func (h *Handler) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	user, _ := h.Auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	pageView.NotFoundPage("", h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path, user).Render(r.Context(), w)
}

func (h *Handler) HandleForbidden(w http.ResponseWriter, r *http.Request) {
	user, _ := h.Auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	pageView.ForbidenPage("", h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path, user).Render(r.Context(), w)
}

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	user, _ := h.Auth.GetAuthenticateUserFromRequest(r, h.Slog)
	pageNav := handlers.GetPageNavDefault(r)

	h.Error = layout.PopinMessages{
		Title:    r.URL.Query().Get("errorTitle"),
		Messages: strings.Split(r.URL.Query().Get("errorMessages"), ","),
	}

	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
		pageNav = h.GetPageNavCustom(r, user, matchModel.Match{})
		pageView.HomePage(user, h.Error, pageNav, h.Languages, r.URL.Path).Render(r.Context(), w)
		return
	}
	pageView.HomePage(userModel.User{}, h.Error, pageNav, h.Languages, r.URL.Path).Render(r.Context(), w)
}

func (h *Handler) HandleCGU(w http.ResponseWriter, r *http.Request) {
	user, _ := h.Auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	legal.CGU(h.Error, h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path).Render(r.Context(), w)
}

func (h *Handler) HandlePrivacy(w http.ResponseWriter, r *http.Request) {
	user, _ := h.Auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if user.Id != 0 {
		h.Slog = h.Slog.With("userId", user.Id)
	}
	legal.Privacy(h.Error, h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path).Render(r.Context(), w)
}
