package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/helper"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/shared/components"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/gorilla/mux"
)

var MaxPlayerByTeam = 8

func (h *Handler) HandlersCreatePlayer(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, _ := model.GetUserByOauth2Id(userOauth2.UserID)
	vars := mux.Vars(r)

	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTeam, _ := strconv.Atoi(r.FormValue("idTeam"))

	if r.FormValue("name") == "" {
		http.Error(w, "Player need a name", http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idTeam"), 10, 64); err != nil {
		http.Error(w, "Player need a color", http.StatusBadRequest)
		return
	}

	canCreatePlayerInTeam, err := model.NumberPlayerInTeamByTeamId(idTeam)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if canCreatePlayerInTeam == MaxPlayerByTeam {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMessages(components.Error{
			Title:    i18n.T(r.Context(), "global.error") + " " + r.FormValue("name"),
			Messages: []string{i18n.T(r.Context(), "page.match.max-player-by-team")},
		}).Render(r.Context(), w)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idClass"), 10, 64); err != nil {
		http.Error(w, "Player need a class", http.StatusBadRequest)
		return
	}

	match, err := model.GetMatch(idMatch)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if match.Round > 0 {
		http.Error(w, "Match already started", http.StatusBadRequest)
		return
	}

	idClass, _ := strconv.Atoi(r.FormValue("idClass"))
	idPlayer, err := model.CreatePlayer(model.PlayerCreate{
		Name:    r.FormValue("name"),
		IdTeam:  idTeam,
		IdClass: idClass,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := model.GetPlayer(user.IdLanguage, idPlayer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.TeamPlayer(player, match).Render(r.Context(), w)
}

func (h *Handler) HandlersDeletePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idPlayer, _ := strconv.Atoi(vars["idPlayer"])

	err := model.DeleteMatchPlayersSpellsByPlayer(idPlayer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.DeletePlayer(idPlayer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}

func (h *Handler) HandlersPlayerLanguage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userOauth2, _ := h.auth.GetSessionUser(r)
	user, _ := model.GetUserByOauth2Id(userOauth2.UserID)

	code := vars["code"]
	idLanguage := helper.SupportedLanguages[code]

	userUpdate := model.UserUpdate{
		IdLanguage: &idLanguage,
	}

	err := model.UpdateUser(user.Id, userUpdate)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Refresh", "true")

}
