package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/helper"
	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/gorilla/mux"
)

var NumberOfMatchMax = 50

func (h *Handler) HandlersListMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	matchs, err := model.GetMatchsByIdUser(user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page.MatchListPage(user, h.error, GetPageNavDefault(r), h.languages, r.URL.Path, matchs).Render(r.Context(), w)
}

func (h *Handler) HandlersCreateMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	err := r.ParseForm()
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))

	numberOfMatch, err := model.GetNumberOfMatchByUserId(user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if name == "" {
		h.Slog.Error("Match need a name", "userId", user.Id, "name", name)
		http.Error(w, "Match need a name", http.StatusBadRequest)
		return
	}

	if numberOfMatch >= NumberOfMatchMax {
		RenderComponentWarning(
			i18n.T(r.Context(), "global.error")+" "+name,
			[]string{i18n.T(r.Context(), "page.match-list.max-match")},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Info("Max number of match", "matchName", name)
		return
	}

	matchId, err := model.CreateMatch(model.MatchCreate{
		Name:   name,
		IdUser: user.Id,
	})

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

	_, err = model.CreateTeam(model.TeamCreate{
		Name:        "Team red",
		IdColorTeam: 1,
		IdMatch:     matchId,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = model.CreateTeam(model.TeamCreate{
		Name:        "Team blue",
		IdColorTeam: 2,
		IdMatch:     matchId,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.Match(match).Render(r.Context(), w)
}

func (h *Handler) HandlersDeleteMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

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
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

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

	page.TeamPlayerListPage(user, h.error, h.GetPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path, match, teams, classes, players).Render(r.Context(), w)
}

func (h *Handler) HandlersMatchUnAutorized(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])

	match, err := model.GetMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.TeamPlayerListPageUnAutorized(user, h.GetPageNavCustom(r, user, model.Match{}), h.languages, r.URL.Path, match).Render(r.Context(), w)
}

func (h *Handler) HandlerStartMatchPage(w http.ResponseWriter, r *http.Request) {
	var idClassGlobal = 13

	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])
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

	if len(players) == 0 {
		http.Error(w, "Match need players", http.StatusBadRequest)
		return
	}

	if match.Round == 0 {
		classeIds := []int{idClassGlobal}
		idSpellToExclude := helper.MasteryIdSpells
		for _, player := range players {
			classeIds = append(classeIds, player.Class.Id)
		}
		spells, err := model.GetSpellsByIdClass(user.IdLanguage, classeIds, idSpellToExclude)
		if err != nil {
			h.Slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if match.MultipleMasteryEnabled == 1 {
			idSpellToExclude = []int{}
		}

		globalSpells, err := model.GetSpellsByIdClass(user.IdLanguage, []int{idClassGlobal}, idSpellToExclude)
		if err != nil {
			h.Slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var round = 1
		err = model.UpdateMatch(matchId, model.MatchUpdate{Round: &round})
		if err != nil {
			h.Slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var matchPlayerSpells []model.MatchPlayerSpellCreate

		for _, player := range players {
			for _, spell := range spells {
				if spell.IdClass == player.Class.Id {
					matchPlayerSpells = append(matchPlayerSpells, model.MatchPlayerSpellCreate{
						MatchId:             matchId,
						PlayerId:            player.Id,
						SpellId:             spell.Id,
						RoundBeforeRecovery: 0,
					})
				}
			}
			for _, globalSpell := range globalSpells {
				matchPlayerSpells = append(matchPlayerSpells, model.MatchPlayerSpellCreate{
					MatchId:             matchId,
					PlayerId:            player.Id,
					SpellId:             globalSpell.Id,
					RoundBeforeRecovery: 0,
				})
			}
		}

		err = model.CreateMatchPlayersSpells(matchPlayerSpells)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	userConfigurationFavoriteSpells, err := model.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayer, err := model.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.StartMatchPage(user, h.error, h.GetPageNavCustom(r, user, match), h.languages, r.URL.Path, match, players, spellsPlayer, false).Render(r.Context(), w)
}

func (h *Handler) HandlerResetMatchPage(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	err := model.ResetMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", matchId))
}

func (h *Handler) HandlerToggleMatchMastery(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])
	multipleMasteryEnabled, _ := strconv.Atoi(vars["toggleBool"])

	err := model.UpdateMatch(matchId, model.MatchUpdate{
		MultipleMasteryEnabled: &multipleMasteryEnabled,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	page.MatchPageTable(user, h.error, h.GetPageNavCustom(r, user, match), h.languages, r.URL.Path, match, players, spellsPlayers, false).Render(r.Context(), w)
}

func (h *Handler) HandlerMatchNextRound(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	err := model.IncreaseMatchRound(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.DecreasePlayersSpellsRoundBeforeRecoveryByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	page.PlayerTable(match, players, spellsPlayers, false).Render(r.Context(), w)
}

func (h *Handler) HandlerUsePlayerSpell(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := model.UsePlayerSpellByIdPlayerSpell(idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellPlayer, err := model.GetSpellPlayerByIdSpellsPlayers(user.IdLanguage, idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.Spell(spellPlayer, false).Render(r.Context(), w)
}

func (h *Handler) HandlerRemoveRoundRecoveryPlayerSpell(w http.ResponseWriter, r *http.Request) {
	user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.Slog)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := model.RemoveRoundRecoverySpellByIdPlayerSpell(idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellPlayer, err := model.GetSpellPlayerByIdSpellsPlayers(user.IdLanguage, idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.Spell(spellPlayer, false).Render(r.Context(), w)
}
