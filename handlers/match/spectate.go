package match

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rousseau-romain/round-timing/config"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/pkg/notify"
	"github.com/rousseau-romain/round-timing/service/auth"
	pageMatch "github.com/rousseau-romain/round-timing/views/page/match"
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

func (h *Handler) HandleSpectateMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	matchFromUser, err := matchModel.GetLastMatchByUserId(user.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	match, err := matchModel.GetMatch(matchId)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(user.IdLanguage, matchId)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurationFavoriteSpells, err := userModel.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := matchModel.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.SpectateMatchPage(user, h.Error, h.GetPageNavCustom(r, user, matchFromUser), h.Languages, r.URL.Path, match, players, spellsPlayers, true).Render(r.Context(), w)
}

func (h *Handler) HandleMatchTableLive(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	ch := notify.Subscribe(matchId)
	defer notify.Unsubscribe(matchId, ch)

	// Detect client disconnect via ReadMessage.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-ch:
			match, err := matchModel.GetMatch(matchId)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			players, err := matchModel.GetPlayersByIdMatch(user.IdLanguage, matchId)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			userConfigurationFavoriteSpells, err := userModel.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			spellsPlayers, err := matchModel.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			message, _ := templ.ToGoHTML(r.Context(), pageMatch.PlayerTable(match, players, spellsPlayers, true))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				logger.Error("Error writing to WebSocket: " + err.Error())
				return
			}
		}
	}
}
