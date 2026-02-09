package match

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/model/game"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/pkg/constants"
	"github.com/rousseau-romain/round-timing/pkg/notify"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/views/page"
	pageMatch "github.com/rousseau-romain/round-timing/views/page/match"
)

type Handler struct {
	*handlers.Handler
}

var NumberOfMatchMax = 50

func (h *Handler) HandleListMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	matchs, err := matchModel.GetMatchsByIdUser(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageMatch.MatchListPage(user, h.Error, handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, matchs).Render(r.Context(), w)
}

func (h *Handler) HandleCreateMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	err := r.ParseForm()
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))

	numberOfMatch, err := matchModel.GetNumberOfMatchByUserId(r.Context(), user.Id)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if name == "" {
		logger.Error("Match need a name", "userId", user.Id, "name", name)
		http.Error(w, "Match need a name", http.StatusBadRequest)
		return
	}

	if numberOfMatch >= NumberOfMatchMax {
		handlers.RenderComponentWarning(
			i18n.T(r.Context(), "global.error")+" "+name,
			[]string{i18n.T(r.Context(), "page.match-list.max-match")},
			http.StatusBadRequest, w, r,
		)
		logger.Info("Max number of match", "matchName", name)
		return
	}

	matchId, err := matchModel.CreateMatch(r.Context(), matchModel.MatchCreate{
		Name:   name,
		IdUser: user.Id,
	})

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("match created", "matchId", matchId, "name", name)

	match, err := matchModel.GetMatch(r.Context(), matchId)

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = matchModel.CreateTeam(r.Context(), matchModel.TeamCreate{
		Name:        "Team red",
		IdColorTeam: 1,
		IdMatch:     matchId,
	})

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = matchModel.CreateTeam(r.Context(), matchModel.TeamCreate{
		Name:        "Team blue",
		IdColorTeam: 2,
		IdMatch:     matchId,
	})

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("teams created", "matchId", matchId)

	pageMatch.Match(match).Render(r.Context(), w)
}

func (h *Handler) HandleDeleteMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId := vars["idMatch"]

	id, _ := strconv.Atoi(matchId)

	err := matchModel.DeleteMatchPlayersSpellsByMatchId(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = matchModel.DeletePlayersByMatchId(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = matchModel.DeleteTeamsByMatchId(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = matchModel.DeleteMatch(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("match deleted", "matchId", id)
}

func (h *Handler) HandleMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])

	match, err := matchModel.GetMatch(r.Context(), matchId)
	if err != nil {
		errorMessage := i18n.T(r.Context(), "page.match.errors.match-not-found", i18n.M{"matchId": matchId})
		logger.Error(errorMessage, "matchId", matchId)
		w.WriteHeader(http.StatusNotFound)
		page.NotFoundPage(errorMessage, h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path, user).Render(r.Context(), w)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(r.Context(), user.IdLanguage, matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := matchModel.GetTeamsByIdMatch(r.Context(), matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := game.GetClasses(r.Context(), user.IdLanguage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.TeamPlayerListPage(user, h.Error, h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path, match, teams, classes, players).Render(r.Context(), w)
}

func (h *Handler) HandleStartMatch(w http.ResponseWriter, r *http.Request) {
	var idClassGlobal = 13

	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])
	match, err := matchModel.GetMatch(r.Context(), matchId)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(r.Context(), user.IdLanguage, matchId)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(players) == 0 {
		http.Error(w, "Match need players", http.StatusBadRequest)
		return
	}

	if match.Round == 0 {
		classeIds := []int{idClassGlobal}
		idSpellToExclude := constants.MasteryIdSpells
		for _, player := range players {
			classeIds = append(classeIds, player.Class.Id)
		}
		spells, err := game.GetSpellsByIdClass(r.Context(), user.IdLanguage, classeIds, idSpellToExclude)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if match.MultipleMasteryEnabled == 1 {
			idSpellToExclude = []int{}
		}

		globalSpells, err := game.GetSpellsByIdClass(r.Context(), user.IdLanguage, []int{idClassGlobal}, idSpellToExclude)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var round = 1
		err = matchModel.UpdateMatch(r.Context(), matchId, matchModel.MatchUpdate{Round: &round})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var matchPlayerSpells []matchModel.MatchPlayerSpellCreate

		spellsByClass := make(map[int][]game.Spell)
		for _, spell := range spells {
			spellsByClass[spell.IdClass] = append(spellsByClass[spell.IdClass], spell)
		}

		for _, player := range players {
			for _, spell := range spellsByClass[player.Class.Id] {
				matchPlayerSpells = append(matchPlayerSpells, matchModel.MatchPlayerSpellCreate{
					MatchId:             matchId,
					PlayerId:            player.Id,
					SpellId:             spell.Id,
					RoundBeforeRecovery: 0,
				})
			}
			for _, globalSpell := range globalSpells {
				matchPlayerSpells = append(matchPlayerSpells, matchModel.MatchPlayerSpellCreate{
					MatchId:             matchId,
					PlayerId:            player.Id,
					SpellId:             globalSpell.Id,
					RoundBeforeRecovery: 0,
				})
			}
		}

		err = matchModel.CreateMatchPlayersSpells(r.Context(), matchPlayerSpells)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("match player spells created", "matchId", matchId, "count", len(matchPlayerSpells))
	}

	logger.Info("match started", "matchId", matchId)

	userConfigurationFavoriteSpells, err := userModel.GetConfigurationByKeyAndIdUser(r.Context(), user.IdLanguage, user.Id, constants.ConfigKeyHideNonFavoriteSpells)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayer, err := matchModel.GetSpellsPlayersByIdMatch(r.Context(), user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.Value == "true")

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notify.Notify(matchId)
	pageMatch.StartMatchPage(user, h.Error, h.GetPageNavCustom(r, user, match), h.Languages, r.URL.Path, match, players, spellsPlayer, false).Render(r.Context(), w)
}

func (h *Handler) HandleResetMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	err := matchModel.ResetMatch(r.Context(), matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("match reset", "matchId", matchId)

	notify.Notify(matchId)
	w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", matchId))
}

func (h *Handler) HandleToggleMatchMastery(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])
	multipleMasteryEnabled, _ := strconv.Atoi(vars["toggleBool"])

	err := matchModel.UpdateMatch(r.Context(), matchId, matchModel.MatchUpdate{
		MultipleMasteryEnabled: &multipleMasteryEnabled,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("match mastery toggled", "matchId", matchId, "multipleMasteryEnabled", multipleMasteryEnabled)

	match, err := matchModel.GetMatch(r.Context(), matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(r.Context(), user.IdLanguage, matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurationFavoriteSpells, err := userModel.GetConfigurationByKeyAndIdUser(r.Context(), user.IdLanguage, user.Id, constants.ConfigKeyHideNonFavoriteSpells)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := matchModel.GetSpellsPlayersByIdMatch(r.Context(), user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.Value == "true")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notify.Notify(matchId)
	pageMatch.MatchPageTable(user, h.Error, h.GetPageNavCustom(r, user, match), h.Languages, r.URL.Path, match, players, spellsPlayers, false).Render(r.Context(), w)
}

func (h *Handler) HandleMatchNextRound(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	err := matchModel.IncreaseMatchRound(r.Context(), matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("match next round", "matchId", matchId)

	err = matchModel.DecreasePlayersSpellsRoundBeforeRecoveryByIdMatch(r.Context(), matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m, err := matchModel.GetMatch(r.Context(), matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(r.Context(), user.IdLanguage, matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurationFavoriteSpells, err := userModel.GetConfigurationByKeyAndIdUser(r.Context(), user.IdLanguage, user.Id, constants.ConfigKeyHideNonFavoriteSpells)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := matchModel.GetSpellsPlayersByIdMatch(r.Context(), user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.Value == "true")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notify.Notify(matchId)
	pageMatch.PlayerTable(m, players, spellsPlayers, false).Render(r.Context(), w)
}

func (h *Handler) HandleUsePlayerSpell(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := matchModel.UsePlayerSpellByIdPlayerSpell(r.Context(), idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("spell used", "matchId", matchId, "playerSpellId", idPlayerSpell)

	spellPlayer, err := matchModel.GetSpellPlayerByIdSpellsPlayers(r.Context(), user.IdLanguage, idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notify.Notify(matchId)
	pageMatch.Spell(spellPlayer, false).Render(r.Context(), w)
}

func (h *Handler) HandleRemoveRoundRecoveryPlayerSpell(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := matchModel.RemoveRoundRecoverySpellByIdPlayerSpell(r.Context(), idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("spell recovery removed", "matchId", matchId, "playerSpellId", idPlayerSpell)

	spellPlayer, err := matchModel.GetSpellPlayerByIdSpellsPlayers(r.Context(), user.IdLanguage, idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notify.Notify(matchId)
	pageMatch.Spell(spellPlayer, false).Render(r.Context(), w)
}
