package profile

import (
	"net/http"

	"github.com/gorilla/mux"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/pkg/lang"
	"github.com/rousseau-romain/round-timing/service/auth"
)

func (h *Handler) HandlePlayerLanguage(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	code := vars["code"]
	idLanguage, ok := lang.SupportedLanguages[code]
	if !ok {
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	userUpdate := userModel.UserUpdate{
		IdLanguage: &idLanguage,
	}

	err := userModel.UpdateUser(r.Context(), user.Id, userUpdate)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("language changed", "code", code)

	w.Header().Set("HX-Refresh", "true")

}
