package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/model"
	pageAdmin "github.com/rousseau-romain/round-timing/views/page/admin"
)

func (h *Handler) HandlersListUser(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)

	users, err := model.GetUsers()
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pageAdmin.UserListPage(user, h.error, GetPageNavDefault(r), h.languages, r.URL.Path, user, users).Render(r.Context(), w)
}

func (h *Handler) HandlersUserEnabled(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idUser, _ := strconv.Atoi(vars["idUser"])

	toggleEnabled, _ := strconv.ParseBool(vars["toggleEnabled"])

	err := model.UpdateUser(idUser, model.UserUpdate{
		Enabled: &toggleEnabled,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := model.GetUserById(idUser)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageAdmin.UserEnabled(user).Render(r.Context(), w)
}
