package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/model/game"
	"github.com/rousseau-romain/round-timing/views/page"
)

func (h *Handler) HandlersToggleSpellFavorite(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idSpell, _ := strconv.Atoi(vars["idSpell"])

	err := game.ToggleIsFavoriteSpell(user.Id, idSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellFavorite, err := game.GetFavoriteSpellByIdUserAndIdSpell(user.IdLanguage, user.Id, idSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.SpellFavorite(spellFavorite).Render(r.Context(), w)
}
