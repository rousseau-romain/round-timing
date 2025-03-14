package handlers

import (
	"fmt"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/google/uuid"
)

func (h *Handler) HandlersProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Slog = h.Slog.With("userId", user.Id)

	idUserShares, err := model.GetUsersSpectateByIdUser(user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := model.GetClasses(user.IdLanguage)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spells, err := model.GetFavoriteSpellsByIdUser(user.IdLanguage, user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurations, err := model.GetAllConfigurationByIdUser(user.IdLanguage, user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.ProfilePage(user, h.error, h.GetPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path, idUserShares, classes, spells, userConfigurations).Render(r.Context(), w)
}

func (h *Handler) HandlersProfileToggleUserConfiguration(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idConfiguration, _ := strconv.Atoi(vars["idConfiguration"])

	err := model.ToggleUserConfiguration(user.Id, idConfiguration)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfiguration, err := model.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, idConfiguration)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.UserConfiguration(userConfiguration).Render(r.Context(), w)
}

func (h *Handler) HandlersProfileAddSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		RenderComponentError(
			"User spectate need a valid id",
			[]string{fmt.Sprintf("User spectate need a valid id not (%s)", r.FormValue("idUserShare"))},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Error("User spectate need a valid id", "userSpectateId", r.FormValue("idUserShare"))
		return
	}

	userSpectateExist, err := model.UserExistsByIdShare(r.FormValue("idUserShare"))

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !userSpectateExist {
		RenderComponentError(
			"User spectate does not exist",
			[]string{"User spectate does not exist"},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Error("User spectate does not exist", "userSpectateId", r.FormValue("idUserShare"))

		return
	}

	IsAlreadyUsersSpectate, err := model.IsUsersSpectateByIdUser(user.Id, r.FormValue("idUserShare"))

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if IsAlreadyUsersSpectate {
		RenderComponentError(
			"User spectate already exist",
			[]string{"User spectate already exist"},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Error("User spectate already exist", "userSpectateId", r.FormValue("idUserShare"))
		return
	}

	_, err = model.CreateUserSpectate(model.UserSpectateCreate{
		IdUser:      user.Id,
		IdUserShare: r.FormValue("idUserShare"),
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.UserSpectate(r.FormValue("idUserShare")).Render(r.Context(), w)
}

func (h *Handler) HandlersProfileDeleteSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		h.Slog.Error("User spectate need a id", "error", err)
		http.Error(w, "User spectate need a id", http.StatusBadRequest)
		return
	}

	if err := model.DeleteUserSpectate(user.Id, r.FormValue("idUserShare")); err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
