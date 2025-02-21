package handlers

import (
	"fmt"

	"log"
	"net/http"

	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/google/uuid"
)

func (h *Handler) HandlersProfile(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, err := model.GetUserByOauth2Id(userOauth2.UserID)
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

	page.ProfilePage(userOauth2, user, h.error, getPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path, idUserShares).Render(r.Context(), w)
}

func (h *Handler) HandlersProfileAddSpectate(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, err := model.GetUserByOauth2Id(userOauth2.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := uuid.Validate(r.FormValue("idUserShare")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		RenderComponentError("User spectate need a valid id", []string{"User spectate need a valid id"}, w, r)
		log.Printf("%s", fmt.Sprintf("User spectate need a valid idShared not (%s)", r.FormValue("idUserShare")))
		return
	}

	userSpectateExist, err := model.UserExistsByIdShare(r.FormValue("idUserShare"))

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !userSpectateExist {
		w.WriteHeader(http.StatusBadRequest)
		RenderComponentError("User spectate does not exist", []string{"User spectate does not exist"}, w, r)
		log.Println("User spectate does not exist")
		return
	}

	IsAlreadyUsersSpectate, err := model.IsUsersSpectateByIdUser(user.Id, r.FormValue("idUserShare"))

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if IsAlreadyUsersSpectate {
		w.WriteHeader(http.StatusBadRequest)
		RenderComponentError("User spectate already exist", []string{"User spectate already exist"}, w, r)
		log.Println("User spectate already exist")
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
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, err := model.GetUserByOauth2Id(userOauth2.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := uuid.Validate(r.FormValue("idUserShare")); err == nil {
		log.Println("User spectate need a id")
		http.Error(w, "User spectate need a id", http.StatusBadRequest)
		return
	}

	err = model.DeleteUserSpectate(user.Id, r.FormValue("idUserShare"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
