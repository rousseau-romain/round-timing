package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/model/game"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/views/page"
)

func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Slog = h.Slog.With("userId", user.Id)

	idUserShares, err := userModel.GetUsersSpectateByIdUser(user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := game.GetClasses(user.IdLanguage)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spells, err := game.GetFavoriteSpellsByIdUser(user.IdLanguage, user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurations, err := userModel.GetAllConfigurationByIdUser(user.IdLanguage, user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.ProfilePage(user, h.error, h.GetPageNavCustom(r, user, matchModel.Match{}), h.languages, r.URL.Path, idUserShares, classes, spells, userConfigurations).Render(r.Context(), w)
}

func (h *Handler) HandleProfileToggleUserConfiguration(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idConfiguration, _ := strconv.Atoi(vars["idConfiguration"])

	err := userModel.ToggleUserConfiguration(user.Id, idConfiguration)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfiguration, err := userModel.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, idConfiguration)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.UserConfiguration(userConfiguration).Render(r.Context(), w)
}

func (h *Handler) HandleProfileAddSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		RenderComponentError(
			i18n.T(r.Context(), "page.profile.errors.user-spectate.not-valid", i18n.M{"userSpectateId": r.FormValue("idUserShare")}),
			[]string{""},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Error("User spectate need a valid id", "userSpectateId", r.FormValue("idUserShare"))
		return
	}

	userSpectateExist, err := userModel.UserExistsByIdShare(r.FormValue("idUserShare"))

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !userSpectateExist {
		RenderComponentError(
			i18n.T(r.Context(), "page.profile.errors.user-spectate.does-not-exist", i18n.M{"userSpectateId": r.FormValue("idUserShare")}),
			[]string{""},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Error("User spectate does not exist", "userSpectateId", r.FormValue("idUserShare"))

		return
	}

	IsAlreadyUsersSpectate, err := userModel.IsUsersSpectateByIdUser(user.Id, r.FormValue("idUserShare"))

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if IsAlreadyUsersSpectate {
		RenderComponentError(
			i18n.T(r.Context(), "page.profile.errors.user-spectate.already-exist", i18n.M{"userSpectateId": r.FormValue("idUserShare")}),
			[]string{""},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Error("User spectate already exist", "userSpectateId", r.FormValue("idUserShare"))
		return
	}

	_, err = userModel.CreateUserSpectate(userModel.UserSpectateCreate{
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

func (h *Handler) HandleProfileDeleteSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		h.Slog.Error("User spectate need a id", "error", err)
		http.Error(w, "User spectate need a id", http.StatusBadRequest)
		return
	}

	if err := userModel.DeleteUserSpectate(user.Id, r.FormValue("idUserShare")); err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
