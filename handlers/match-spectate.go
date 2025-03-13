package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rousseau-romain/round-timing/config"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			config.PUBLIC_HOST_PORT: true,
		}
		return allowedOrigins[origin]
	},
}

var clients = make(map[string]*websocket.Conn) // Maps userID to WebSocket connection

func (h *Handler) HandlerSpectateMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	matchFromUser, err := model.GetLastMatchByUserId(user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	match, err := model.GetMatch(matchId)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := model.GetPlayersByIdMatch(user.IdLanguage, matchId)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurationFavoriteSpells, err := model.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := model.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.SpectateMatchPage(user, h.error, h.GetPageNavCustom(r, user, matchFromUser), h.languages, r.URL.Path, match, players, spellsPlayers, true).Render(r.Context(), w)
}

func (h *Handler) HandlerMatchTableLive(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])
	userMatchUniqueString := fmt.Sprintf("%d-%d", user.Id, matchId)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	clients[userMatchUniqueString] = conn

	defer func() {
		delete(clients, userMatchUniqueString)
		conn.Close()
	}()

	for {
		time.Sleep(2 * time.Second)

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

		userConfigurationFavoriteSpells, err := model.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
		if err != nil {
			h.Slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		spellsPlayers, err := model.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		message, _ := templ.ToGoHTML(r.Context(), page.PlayerTable(match, players, spellsPlayers, true))
		for client := range clients {
			c := clients[userMatchUniqueString]
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Error writing to WebSocket:", err)
				c.Close()
				delete(clients, client)
			}
		}
	}
}
