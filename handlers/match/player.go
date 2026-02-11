package match

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/handlers"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/pkg/notify"
	"github.com/rousseau-romain/round-timing/service/auth"
	pageMatch "github.com/rousseau-romain/round-timing/views/page/match"
)

var MaxPlayerByTeam = 8

func (h *Handler) HandleUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])
	name := strings.TrimSpace(r.FormValue("name"))
	idPlayer, _ := strconv.Atoi(vars["idPlayer"])
	if name == "" {
		logger.Info("Player need a name", "idPlayer", idPlayer, "name", name)
		http.Error(w, "Player need a name", http.StatusBadRequest)
		return
	}

	err := matchModel.UpdatePlayer(r.Context(), idPlayer, matchModel.PlayerUpdate{
		Name: &name,
	})
	if err != nil {
		httpError.InternalError(w)
		logger.Error(err.Error())
		return
	}

	logger.Info("player updated", "playerId", idPlayer, "name", name)

	notify.Notify(matchId)
}

func (h *Handler) HandleCreatePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTeam, _ := strconv.Atoi(r.FormValue("idTeam"))
	name := strings.TrimSpace(r.FormValue("name"))

	if name == "" {
		logger.Info("Player need a name", "idPlayer", user.Id, "name", name)
		http.Error(w, "Player need a name", http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idTeam"), 10, 64); err != nil {
		logger.Info("Player need a color", "idPlayer", user.Id, "name", name)
		http.Error(w, "Player need a color", http.StatusBadRequest)
		return
	}

	canCreatePlayerInTeam, err := matchModel.NumberPlayerInTeamByTeamId(r.Context(), idTeam)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
	}

	if canCreatePlayerInTeam == MaxPlayerByTeam {
		handlers.RenderComponentWarning(
			i18n.T(r.Context(), "global.error")+" "+name,
			[]string{i18n.T(r.Context(), "page.match.max-player-by-team")},
			http.StatusBadRequest,
			w, r,
		)
		logger.Info("Max player in team", "teamId", idTeam, "playerName", name)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idClass"), 10, 64); err != nil {
		logger.Info("Player need a class", "idPlayer", user.Id, "name", name)
		http.Error(w, "Player need a class", http.StatusBadRequest)
		return
	}

	match, err := matchModel.GetMatch(r.Context(), idMatch)

	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	if match.Round > 0 {
		logger.Info("Match already started", "idPlayer", user.Id, "name", name)
		http.Error(w, "Match already started", http.StatusBadRequest)
		return
	}

	idClass, _ := strconv.Atoi(r.FormValue("idClass"))
	idPlayer, err := matchModel.CreatePlayer(r.Context(), matchModel.PlayerCreate{
		Name:    name,
		IdTeam:  idTeam,
		IdClass: idClass,
	})

	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("player created", "playerId", idPlayer, "name", name, "teamId", idTeam)

	player, err := matchModel.GetPlayer(r.Context(), user.IdLanguage, idPlayer)

	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	notify.Notify(idMatch)
	pageMatch.TeamPlayer(player, match).Render(r.Context(), w)
}

func (h *Handler) HandleDeletePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])
	idPlayer, _ := strconv.Atoi(vars["idPlayer"])

	err := matchModel.DeleteMatchPlayersSpellsByPlayer(r.Context(), idPlayer)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	err = matchModel.DeletePlayer(r.Context(), idPlayer)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("player deleted", "playerId", idPlayer, "matchId", matchId)

	notify.Notify(matchId)
}
