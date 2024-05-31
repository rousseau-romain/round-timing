package handlers

import (
	"log"
	"net/http"
	"round-timing/model"
	"round-timing/views/page"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) HandlersCreatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idMatch, _ := strconv.Atoi(vars["idMatch"])

	if r.FormValue("name") == "" {
		http.Error(w, "Team need a name", http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idTeam"), 10, 64); err != nil {
		http.Error(w, "Team need a color", http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseInt(r.FormValue("idClass"), 10, 64); err != nil {
		http.Error(w, "Team need a class", http.StatusBadRequest)
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

	idTeam, _ := strconv.Atoi(r.FormValue("idTeam"))
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

	player, err := model.GetPlayer(idPlayer)

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
