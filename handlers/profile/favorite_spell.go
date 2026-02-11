package profile

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/model/game"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/views/page"
)

func (h *Handler) HandleToggleSpellFavorite(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idSpell, _ := strconv.Atoi(vars["idSpell"])

	err := game.ToggleIsFavoriteSpell(r.Context(), user.Id, idSpell)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("spell favorite toggled", "spellId", idSpell)

	spellFavorite, err := game.GetFavoriteSpellByIdUserAndIdSpell(r.Context(), user.IdLanguage, user.Id, idSpell)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	page.SpellFavorite(spellFavorite).Render(r.Context(), w)
}
