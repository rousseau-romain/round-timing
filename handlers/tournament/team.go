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

func (h *Handler) HandleAddTeam(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		handlers.RenderComponentError("Team needs a name", []string{"Team needs a name"}, http.StatusBadRequest, w, r)
		return
	}

	id, err := tournamentModel.CreateTeam(r.Context(), tournamentModel.TeamCreate{
		IdUser: user.Id,
		Name:   name,
	})
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("team added", "teamId", id, "name", name)

	t, err := tournamentModel.GetTeam(r.Context(), id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.TeamCard(t).Render(r.Context(), w)
}

func (h *Handler) HandleDeleteTeam(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTeam, _ := strconv.Atoi(vars["idTeam"])

	if err := tournamentModel.DeleteTeam(r.Context(), idTeam, user.Id); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("team deleted", "teamId", idTeam)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}
