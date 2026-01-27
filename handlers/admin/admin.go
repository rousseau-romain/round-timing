package admin

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/service/auth"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	pageAdmin "github.com/rousseau-romain/round-timing/views/page/admin"
)

type Handler struct {
	*handlers.Handler
}

func (h *Handler) HandleListUser(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	users, err := userModel.GetUsers()
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pageAdmin.UserListPage(user, h.Error, handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, user, users).Render(r.Context(), w)
}

func (h *Handler) HandleUserEnabled(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idUser, _ := strconv.Atoi(vars["idUser"])

	toggleEnabled, _ := strconv.ParseBool(vars["toggleEnabled"])

	err := userModel.UpdateUser(idUser, userModel.UserUpdate{
		Enabled: &toggleEnabled,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userEnabled, err := userModel.GetUserById(idUser)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageAdmin.UserEnabled(userEnabled).Render(r.Context(), w)
}
