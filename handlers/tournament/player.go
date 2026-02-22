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

func (h *Handler) HandleAddPlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		handlers.RenderComponentError("Player needs a name", []string{"Player needs a name"}, http.StatusBadRequest, w, r)
		return
	}

	id, err := tournamentModel.CreatePlayer(r.Context(), tournamentModel.PlayerCreate{
		IdUser: user.Id,
		Name:   name,
	})
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("player added", "playerId", id, "name", name)

	p, err := tournamentModel.GetPlayer(r.Context(), id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.PlayerCard(p).Render(r.Context(), w)
}

func (h *Handler) HandleDeletePlayer(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idPlayer, _ := strconv.Atoi(vars["idPlayer"])

	if err := tournamentModel.DeletePlayer(r.Context(), idPlayer, user.Id); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("player deleted", "playerId", idPlayer)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}
