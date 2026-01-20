package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/helper"
	"github.com/rousseau-romain/round-timing/model"
	pageMatch "github.com/rousseau-romain/round-timing/views/page/match"

	"github.com/gorilla/mux"
)

var MaxPlayerByTeam = 8

func (h *Handler) HandlersUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	name := strings.TrimSpace(r.FormValue("name"))
	idPlayer, _ := strconv.Atoi(vars["idPlayer"])
	if name == "" {
		h.Slog.Info("Player need a name", "idPlayer", idPlayer, "name", name)
		http.Error(w, "Player need a name", http.StatusBadRequest)
		return
	}

	err := model.UpdatePlayer(idPlayer, model.PlayerUpdate{
		Name: &name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Slog.Error(err.Error())
		return
	}

}

func (h *Handler) HandlersCreatePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTeam, _ := strconv.Atoi(r.FormValue("idTeam"))
	name := strings.TrimSpace(r.FormValue("name"))

	if name == "" {
		h.Slog.Info("Player need a name", "idPlayer", user.Id, "name", name)
		http.Error(w, "Player need a name", http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idTeam"), 10, 64); err != nil {
		h.Slog.Info("Player need a color", "idPlayer", user.Id, "name", name)
		http.Error(w, "Player need a color", http.StatusBadRequest)
		return
	}

	canCreatePlayerInTeam, err := model.NumberPlayerInTeamByTeamId(idTeam)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if canCreatePlayerInTeam == MaxPlayerByTeam {
		RenderComponentWarning(
			i18n.T(r.Context(), "global.error")+" "+name,
			[]string{i18n.T(r.Context(), "page.match.max-player-by-team")},
			http.StatusBadRequest,
			w, r,
		)
		h.Slog.Info("Max player in team", "teamId", idTeam, "playerName", name)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idClass"), 10, 64); err != nil {
		h.Slog.Info("Player need a class", "idPlayer", user.Id, "name", name)
		http.Error(w, "Player need a class", http.StatusBadRequest)
		return
	}

	match, err := model.GetMatch(idMatch)

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if match.Round > 0 {
		h.Slog.Info("Match already started", "idPlayer", user.Id, "name", name)
		http.Error(w, "Match already started", http.StatusBadRequest)
		return
	}

	idClass, _ := strconv.Atoi(r.FormValue("idClass"))
	idPlayer, err := model.CreatePlayer(model.PlayerCreate{
		Name:    name,
		IdTeam:  idTeam,
		IdClass: idClass,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := model.GetPlayer(user.IdLanguage, idPlayer)

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.TeamPlayer(player, match).Render(r.Context(), w)
}

func (h *Handler) HandlersDeletePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	idPlayer, _ := strconv.Atoi(vars["idPlayer"])

	err := model.DeleteMatchPlayersSpellsByPlayer(idPlayer)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.DeletePlayer(idPlayer)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandlersPlayerLanguage(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	code := vars["code"]
	idLanguage := helper.SupportedLanguages[code]

	userUpdate := model.UserUpdate{
		IdLanguage: &idLanguage,
	}

	err := model.UpdateUser(user.Id, userUpdate)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Refresh", "true")

}
