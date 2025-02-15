package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/shared/components"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/gorilla/mux"
)

var NumberOfMatchMax = 50

func (h *Handler) HandlersListMatch(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, _ := model.GetUserByOauth2Id(userOauth2.UserID)

	matchs, err := model.GetMatchsByIdUser(user.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page.MatchListPage(userOauth2, user, h.error, GetPageNavDefault(r), h.languages, user, matchs).Render(r.Context(), w)
}

func (h *Handler) HandlersCreateMatch(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, _ := model.GetUserByOauth2Id(userOauth2.UserID)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	numberOfMatch, err := model.GetNumberOfMatchByUserId(user.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.TrimSpace(r.FormValue("name")) == "" {
		log.Println(err)
		http.Error(w, "Match need a name", http.StatusBadRequest)
		return
	}

	if numberOfMatch >= NumberOfMatchMax {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMessages(components.Error{
			Title:    i18n.T(r.Context(), "global.error") + " " + r.FormValue("name"),
			Messages: []string{i18n.T(r.Context(), "page.match-list.max-match")},
		}).Render(r.Context(), w)
		return
	}

	matchId, err := model.CreateMatch(model.MatchCreate{
		Name:   r.FormValue("name"),
		IdUser: user.Id,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	match, err := model.GetMatch(matchId)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = model.CreateTeam(model.TeamCreate{
		Name:        "Team red",
		IdColorTeam: 1,
		IdMatch:     matchId,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = model.CreateTeam(model.TeamCreate{
		Name:        "Team blue",
		IdColorTeam: 2,
		IdMatch:     matchId,
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.Match(match).Render(r.Context(), w)
}

func (h *Handler) HandlersDeleteMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchId := vars["idMatch"]

	id, _ := strconv.Atoi(matchId)

	err := model.DeleteMatchPlayersSpellsByMatchId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.DeletePlayersByMatchId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.DeleteTeamsByMatchId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.DeleteMatch(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandlersMatch(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, _ := model.GetUserByOauth2Id(userOauth2.UserID)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])

	match, err := model.GetMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := model.GetPlayersByIdMatch(user.IdLanguage, matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := model.GetTeamsByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := model.GetClasses(user.IdLanguage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.TeamPlayerListPage(userOauth2, user, h.error, getPageNavCustom(r, user, model.Match{}), h.languages, user, match, teams, classes, players).Render(r.Context(), w)
}
