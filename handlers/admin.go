package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/model"
	pageAdmin "github.com/rousseau-romain/round-timing/views/page/admin"
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

	pageAdmin.UserListPage(userOauth2, user, h.error, GetPageNavDefault(r), h.languages, r.URL.Path, user, users).Render(r.Context(), w)
}

func (h *Handler) HandlersUserEnabled(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idUser, _ := strconv.Atoi(vars["idUser"])

	toggleEnabled, _ := strconv.ParseBool(vars["toggleEnabled"])

	err := model.UpdateUser(idUser, model.UserUpdate{
		Enabled: &toggleEnabled,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := model.GetUserById(idUser)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageAdmin.UserEnabled(user).Render(r.Context(), w)
}
