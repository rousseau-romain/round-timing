package tournament

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	tournamentModel "github.com/rousseau-romain/round-timing/model/tournament"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/service/auth"
	pageTournament "github.com/rousseau-romain/round-timing/views/page/tournament"
)

type Handler struct {
	*handlers.Handler
}

func (h *Handler) HandleListTournament(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	tournaments, err := tournamentModel.GetTournamentsByIdUser(r.Context(), user.Id, false)
	if err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	pageTournament.TournamentListPage(user, h.Error, handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, tournaments).Render(r.Context(), w)
}

func (h *Handler) HandleCreateTournament(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		handlers.RenderComponentError("Tournament needs a name", []string{"Tournament needs a name"}, http.StatusBadRequest, w, r)
		return
	}

	numberPlayerByTeam, err := strconv.Atoi(r.FormValue("number_player_by_team"))
	if err != nil || numberPlayerByTeam < 1 {
		numberPlayerByTeam = 1
	}

	stageType := r.FormValue("stage_type")
	if stageType == "" {
		stageType = "draft"
	}

	tournamentId, err := tournamentModel.CreateTournament(r.Context(), tournamentModel.TournamentCreate{
		Name:               name,
		IdUser:             user.Id,
		NumberPlayerByTeam: numberPlayerByTeam,
		StageType:          stageType,
	})
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("tournament created", "tournamentId", tournamentId, "name", name)

	t, err := tournamentModel.GetTournament(r.Context(), tournamentId)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.TournamentRow(t).Render(r.Context(), w)
}

func (h *Handler) HandleViewTournament(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	teams, err := tournamentModel.GetTeamsByUser(r.Context(), user.Id)
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

	teamPlayers, err := tournamentModel.GetTeamPlayersByUser(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.TournamentViewPage(user, h.Error, handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, t, teams, players, teamPlayers).Render(r.Context(), w)
}

func (h *Handler) HandleViewTournamentMatchs(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.TournamentMatchsPage(user, h.Error, handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, t, matches, availableTeams).Render(r.Context(), w)
}

func (h *Handler) HandleDeleteTournament(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := tournamentModel.DeleteTournamentMatchesByTournament(r.Context(), idTournament); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	if err := tournamentModel.DeleteTournament(r.Context(), idTournament); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("tournament deleted", "tournamentId", idTournament)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}

func (h *Handler) HandleListArchivedTournament(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	tournaments, err := tournamentModel.GetArchivedTournamentsByIdUser(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	pageTournament.ArchivedTournamentSection(tournaments).Render(r.Context(), w)
}

func (h *Handler) HandleHideArchivedSection(w http.ResponseWriter, r *http.Request) {
	pageTournament.ArchivedSectionToggle().Render(r.Context(), w)
}

func (h *Handler) HandleArchiveTournament(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := tournamentModel.ArchiveTournament(r.Context(), idTournament); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("tournament archived", "tournamentId", idTournament)

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	archivedTournaments, err := tournamentModel.GetArchivedTournamentsByIdUser(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.ArchiveTournamentResponse(t, archivedTournaments).Render(r.Context(), w)
}

func (h *Handler) HandleUnarchiveTournament(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := tournamentModel.UnarchiveTournament(r.Context(), idTournament); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("tournament unarchived", "tournamentId", idTournament)

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	archivedTournaments, err := tournamentModel.GetArchivedTournamentsByIdUser(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.UnarchiveTournamentResponse(t, archivedTournaments).Render(r.Context(), w)
}

func (h *Handler) HandleTournamentStatusBadge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		h.Slog.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.StatusBadge(t.Status).Render(r.Context(), w)
}
