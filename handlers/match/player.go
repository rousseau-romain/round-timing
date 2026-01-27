package match

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/service/auth"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	pageMatch "github.com/rousseau-romain/round-timing/views/page/match"
)

var MaxPlayerByTeam = 8

func (h *Handler) HandleUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	name := strings.TrimSpace(r.FormValue("name"))
	idPlayer, _ := strconv.Atoi(vars["idPlayer"])
	if name == "" {
		h.Slog.Info("Player need a name", "idPlayer", idPlayer, "name", name)
		http.Error(w, "Player need a name", http.StatusBadRequest)
		return
	}

	err := matchModel.UpdatePlayer(idPlayer, matchModel.PlayerUpdate{
		Name: &name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Slog.Error(err.Error())
		return
	}

}

func (h *Handler) HandleCreatePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
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

	canCreatePlayerInTeam, err := matchModel.NumberPlayerInTeamByTeamId(idTeam)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if canCreatePlayerInTeam == MaxPlayerByTeam {
		handlers.RenderComponentWarning(
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

	match, err := matchModel.GetMatch(idMatch)

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
	idPlayer, err := matchModel.CreatePlayer(matchModel.PlayerCreate{
		Name:    name,
		IdTeam:  idTeam,
		IdClass: idClass,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := matchModel.GetPlayer(user.IdLanguage, idPlayer)

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.TeamPlayer(player, match).Render(r.Context(), w)
}

func (h *Handler) HandleDeletePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	idPlayer, _ := strconv.Atoi(vars["idPlayer"])

	err := matchModel.DeleteMatchPlayersSpellsByPlayer(idPlayer)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = matchModel.DeletePlayer(idPlayer)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
