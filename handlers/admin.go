package handlers

import (
	"log"
	"net/http"

	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"
)

func (h *Handler) HandlersListUser(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, _ := model.GetUserByOauth2Id(userOauth2.UserID)

	users, err := model.GetUsers()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page.UserListPage(userOauth2, user, h.error, PagesNav, user, users).Render(r.Context(), w)
}
