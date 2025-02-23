package handlers

import (
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"
)

func (h *Handler) HandlersToggleSpellFavorite(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, _ := model.GetUserByOauth2Id(userOauth2.UserID)
	vars := mux.Vars(r)
	idSpell, _ := strconv.Atoi(vars["idSpell"])

	err := model.ToggleIsFavoriteSpell(user.Id, idSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellFavorite, err := model.GetFavoriteSpellByIdUserAndIdSpell(user.IdLanguage, user.Id, idSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.SpellFavorite(spellFavorite).Render(r.Context(), w)
}
