package profile

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/model/game"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/views/components/ui"
	"github.com/rousseau-romain/round-timing/views/page"
)

type Handler struct {
	*handlers.Handler
}

func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromRequest(r)
	if !ok {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}
	logger := h.Slog.With("userId", user.Id)

	idUserShares, err := userModel.GetUsersSpectateByIdUser(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := game.GetClasses(r.Context(), user.IdLanguage)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spells, err := game.GetFavoriteSpellsByIdUser(r.Context(), user.IdLanguage, user.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurations, err := userModel.GetAllConfigurationByIdUser(r.Context(), user.IdLanguage, user.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.ProfilePage(user, h.Error, h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path, idUserShares, classes, spells, userConfigurations).Render(r.Context(), w)
}

func (h *Handler) HandleProfileToggleUserConfiguration(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idConfiguration, _ := strconv.Atoi(vars["idConfiguration"])

	value := r.FormValue("value")
	if value != "" {
		err := userModel.SetUserConfiguration(r.Context(), user.Id, idConfiguration, value)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Info("user configuration set", "configurationId", idConfiguration, "value", value)
	} else {
		err := userModel.ToggleUserConfiguration(r.Context(), user.Id, idConfiguration, user.IdLanguage)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Info("user configuration toggled", "configurationId", idConfiguration)
	}

	userConfiguration, err := userModel.GetConfigurationByIdConfigurationIdUser(r.Context(), user.IdLanguage, user.Id, idConfiguration)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.UserConfiguration(userConfiguration).Render(r.Context(), w)
}

func (h *Handler) HandleToggleContainerExpanded(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idConfiguration, _ := strconv.Atoi(vars["idConfiguration"])

	err := userModel.ToggleUserConfiguration(r.Context(), user.Id, idConfiguration, user.IdLanguage)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("container expanded toggled", "configurationId", idConfiguration)

	configs, err := userModel.GetAllConfigurationByIdUser(r.Context(), user.IdLanguage, user.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	configId, isExpanded := ui.GetContainerExpandedConfig(configs)
	ui.ButtonExpandWidth(configId, isExpanded).Render(r.Context(), w)
}

func (h *Handler) HandleProfileAddSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		handlers.RenderComponentError(
			i18n.T(r.Context(), "page.profile.errors.user-spectate.not-valid", i18n.M{"userSpectateId": r.FormValue("idUserShare")}),
			[]string{""},
			http.StatusBadRequest, w, r,
		)
		logger.Error("User spectate need a valid id", "userSpectateId", r.FormValue("idUserShare"))
		return
	}

	userSpectateExist, err := userModel.UserExistsByIdShare(r.Context(), r.FormValue("idUserShare"))

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !userSpectateExist {
		handlers.RenderComponentError(
			i18n.T(r.Context(), "page.profile.errors.user-spectate.does-not-exist", i18n.M{"userSpectateId": r.FormValue("idUserShare")}),
			[]string{""},
			http.StatusBadRequest, w, r,
		)
		logger.Error("User spectate does not exist", "userSpectateId", r.FormValue("idUserShare"))

		return
	}

	IsAlreadyUsersSpectate, err := userModel.IsUsersSpectateByIdUser(r.Context(), user.Id, r.FormValue("idUserShare"))

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if IsAlreadyUsersSpectate {
		handlers.RenderComponentError(
			i18n.T(r.Context(), "page.profile.errors.user-spectate.already-exist", i18n.M{"userSpectateId": r.FormValue("idUserShare")}),
			[]string{""},
			http.StatusBadRequest, w, r,
		)
		logger.Error("User spectate already exist", "userSpectateId", r.FormValue("idUserShare"))
		return
	}

	_, err = userModel.CreateUserSpectate(r.Context(), userModel.UserSpectateCreate{
		IdUser:      user.Id,
		IdUserShare: r.FormValue("idUserShare"),
	})

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("user spectate created", "userSpectateId", r.FormValue("idUserShare"))

	page.UserSpectate(r.FormValue("idUserShare")).Render(r.Context(), w)
}

func (h *Handler) HandleProfileDeleteSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		logger.Error("User spectate need a id", "error", err)
		http.Error(w, "User spectate need a id", http.StatusBadRequest)
		return
	}

	if err := userModel.DeleteUserSpectate(r.Context(), user.Id, r.FormValue("idUserShare")); err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("user spectate deleted", "userSpectateId", r.FormValue("idUserShare"))
}
