package handlers

import (
	"fmt"
	"strconv"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/google/uuid"
)

func (h *Handler) HandlersProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetAuthenticateUserFromRequest(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idUserShares, err := model.GetUsersSpectateByIdUser(user.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := model.GetClasses(user.IdLanguage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spells, err := model.GetFavoriteSpellsByIdUser(user.IdLanguage, user.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurations, err := model.GetAllConfigurationByIdUser(user.IdLanguage, user.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.ProfilePage(user, h.error, getPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path, idUserShares, classes, spells, userConfigurations).Render(r.Context(), w)
}

func (h *Handler) HandlersProfileToggleUserConfiguration(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r)
	vars := mux.Vars(r)
	idConfiguration, _ := strconv.Atoi(vars["idConfiguration"])

	err := model.ToggleUserConfiguration(user.Id, idConfiguration)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfiguration, err := model.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, idConfiguration)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.UserConfiguration(userConfiguration).Render(r.Context(), w)
}

func (h *Handler) HandlersProfileAddSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		RenderComponentErrorAndLog(
			"User spectate need a valid id",
			[]string{"User spectate need a valid id"},
			[]string{fmt.Sprintf("User spectate need a valid idShared not (%s)", r.FormValue("idUserShare"))},
			http.StatusBadRequest, w, r,
		)
		return
	}

	userSpectateExist, err := model.UserExistsByIdShare(r.FormValue("idUserShare"))

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !userSpectateExist {
		RenderComponentErrorAndLog(
			"User spectate does not exist",
			[]string{"User spectate does not exist"},
			[]string{"User spectate does not exist"},
			http.StatusBadRequest, w, r,
		)
		return
	}

	IsAlreadyUsersSpectate, err := model.IsUsersSpectateByIdUser(user.Id, r.FormValue("idUserShare"))

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if IsAlreadyUsersSpectate {
		RenderComponentErrorAndLog(
			"User spectate already exist",
			[]string{"User spectate already exist"},
			[]string{"User spectate already exist"},
			http.StatusBadRequest, w, r,
		)
		return
	}

	_, err = model.CreateUserSpectate(model.UserSpectateCreate{
		IdUser:      user.Id,
		IdUserShare: r.FormValue("idUserShare"),
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.UserSpectate(r.FormValue("idUserShare")).Render(r.Context(), w)
}

func (h *Handler) HandlersProfileDeleteSpectate(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r)

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		log.Println("User spectate need a id")
		http.Error(w, "User spectate need a id", http.StatusBadRequest)
		return
	}

	if err := model.DeleteUserSpectate(user.Id, r.FormValue("idUserShare")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
