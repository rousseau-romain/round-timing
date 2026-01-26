package profile

import (
	"net/http"

	"github.com/gorilla/mux"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/pkg/lang"
)

func (h *Handler) HandlePlayerLanguage(w http.ResponseWriter, r *http.Request) {
	user, _ := h.Auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	code := vars["code"]
	idLanguage := lang.SupportedLanguages[code]

	userUpdate := userModel.UserUpdate{
		IdLanguage: &idLanguage,
	}

	err := userModel.UpdateUser(user.Id, userUpdate)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Refresh", "true")

}
