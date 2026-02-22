package tournament

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	tournamentModel "github.com/rousseau-romain/round-timing/model/tournament"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/service/auth"
	pageTournament "github.com/rousseau-romain/round-timing/views/page/tournament"
)

func (h *Handler) HandleAddTeamPlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTeam, _ := strconv.Atoi(vars["idTeam"])

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	idPlayer, err := strconv.Atoi(r.FormValue("id_player"))
	if err != nil {
		http.Error(w, "Invalid player", http.StatusBadRequest)
		return
	}

	id, err := tournamentModel.CreateTeamPlayer(r.Context(), tournamentModel.TeamPlayerCreate{
		IdUser:   user.Id,
		IdTeam:   idTeam,
		IdPlayer: idPlayer,
	})
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("player added to team", "teamId", idTeam, "playerId", idPlayer, "teamPlayerId", id)

	teamPlayers, err := tournamentModel.GetTeamPlayersByTeam(r.Context(), idTeam)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	players, err := tournamentModel.GetPlayersByUser(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	team, err := tournamentModel.GetTeam(r.Context(), idTeam)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.TeamCompositionCard(team, teamPlayers, players).Render(r.Context(), w)
}

func (h *Handler) HandleDeleteTeamPlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTeamPlayer, _ := strconv.Atoi(vars["idTeamPlayer"])

	if err := tournamentModel.DeleteTeamPlayer(r.Context(), idTeamPlayer, user.Id); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("team player deleted", "teamPlayerId", idTeamPlayer)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}
