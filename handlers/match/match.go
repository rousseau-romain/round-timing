package match

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/invopop/ctxi18n/i18n"
	"github.com/rousseau-romain/round-timing/handlers"
	"github.com/rousseau-romain/round-timing/service/auth"
	"github.com/rousseau-romain/round-timing/model/game"
	matchModel "github.com/rousseau-romain/round-timing/model/match"
	userModel "github.com/rousseau-romain/round-timing/model/user"
	"github.com/rousseau-romain/round-timing/pkg/constants"
	"github.com/rousseau-romain/round-timing/views/page"
	pageMatch "github.com/rousseau-romain/round-timing/views/page/match"
)

type Handler struct {
	*handlers.Handler
}

var NumberOfMatchMax = 50

func (h *Handler) HandleListMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	matchs, err := matchModel.GetMatchsByIdUser(user.Id)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pageMatch.MatchListPage(user, h.Error, handlers.GetPageNavDefault(r), h.Languages, r.URL.Path, matchs).Render(r.Context(), w)
}

func (h *Handler) HandleCreateMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	err := r.ParseForm()
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))

	numberOfMatch, err := matchModel.GetNumberOfMatchByUserId(user.Id)
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
		handlers.RenderComponentWarning(
			i18n.T(r.Context(), "global.error")+" "+name,
			[]string{i18n.T(r.Context(), "page.match-list.max-match")},
			http.StatusBadRequest, w, r,
		)
		h.Slog.Info("Max number of match", "matchName", name)
		return
	}

	matchId, err := matchModel.CreateMatch(matchModel.MatchCreate{
		Name:   name,
		IdUser: user.Id,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	match, err := matchModel.GetMatch(matchId)

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = matchModel.CreateTeam(matchModel.TeamCreate{
		Name:        "Team red",
		IdColorTeam: 1,
		IdMatch:     matchId,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = matchModel.CreateTeam(matchModel.TeamCreate{
		Name:        "Team blue",
		IdColorTeam: 2,
		IdMatch:     matchId,
	})

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.Match(match).Render(r.Context(), w)
}

func (h *Handler) HandleDeleteMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId := vars["idMatch"]

	id, _ := strconv.Atoi(matchId)

	err := matchModel.DeleteMatchPlayersSpellsByMatchId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = matchModel.DeletePlayersByMatchId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = matchModel.DeleteTeamsByMatchId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = matchModel.DeleteMatch(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])

	match, err := matchModel.GetMatch(matchId)
	if err != nil {
		errorMessage := i18n.T(r.Context(), "page.match.errors.match-not-found", i18n.M{"matchId": matchId})
		h.Slog.Error(errorMessage, "matchId", matchId)
		w.WriteHeader(http.StatusNotFound)
		page.NotFoundPage(errorMessage, h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path, user).Render(r.Context(), w)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(user.IdLanguage, matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := matchModel.GetTeamsByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := game.GetClasses(user.IdLanguage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.TeamPlayerListPage(user, h.Error, h.GetPageNavCustom(r, user, matchModel.Match{}), h.Languages, r.URL.Path, match, teams, classes, players).Render(r.Context(), w)
}

func (h *Handler) HandleStartMatch(w http.ResponseWriter, r *http.Request) {
	var idClassGlobal = 13

	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])
	match, err := matchModel.GetMatch(matchId)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(user.IdLanguage, matchId)
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
		idSpellToExclude := constants.MasteryIdSpells
		for _, player := range players {
			classeIds = append(classeIds, player.Class.Id)
		}
		spells, err := game.GetSpellsByIdClass(user.IdLanguage, classeIds, idSpellToExclude)
		if err != nil {
			h.Slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if match.MultipleMasteryEnabled == 1 {
			idSpellToExclude = []int{}
		}

		globalSpells, err := game.GetSpellsByIdClass(user.IdLanguage, []int{idClassGlobal}, idSpellToExclude)
		if err != nil {
			h.Slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var round = 1
		err = matchModel.UpdateMatch(matchId, matchModel.MatchUpdate{Round: &round})
		if err != nil {
			h.Slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var matchPlayerSpells []matchModel.MatchPlayerSpellCreate

		for _, player := range players {
			for _, spell := range spells {
				if spell.IdClass == player.Class.Id {
					matchPlayerSpells = append(matchPlayerSpells, matchModel.MatchPlayerSpellCreate{
						MatchId:             matchId,
						PlayerId:            player.Id,
						SpellId:             spell.Id,
						RoundBeforeRecovery: 0,
					})
				}
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

		err = matchModel.CreateMatchPlayersSpells(matchPlayerSpells)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	userConfigurationFavoriteSpells, err := userModel.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayer, err := matchModel.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)

	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.StartMatchPage(user, h.Error, h.GetPageNavCustom(r, user, match), h.Languages, r.URL.Path, match, players, spellsPlayer, false).Render(r.Context(), w)
}

func (h *Handler) HandleResetMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	err := matchModel.ResetMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", matchId))
}

func (h *Handler) HandleToggleMatchMastery(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])
	multipleMasteryEnabled, _ := strconv.Atoi(vars["toggleBool"])

	err := matchModel.UpdateMatch(matchId, matchModel.MatchUpdate{
		MultipleMasteryEnabled: &multipleMasteryEnabled,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	match, err := matchModel.GetMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(user.IdLanguage, matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurationFavoriteSpells, err := userModel.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := matchModel.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.MatchPageTable(user, h.Error, h.GetPageNavCustom(r, user, match), h.Languages, r.URL.Path, match, players, spellsPlayers, false).Render(r.Context(), w)
}

func (h *Handler) HandleMatchNextRound(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])

	err := matchModel.IncreaseMatchRound(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = matchModel.DecreasePlayersSpellsRoundBeforeRecoveryByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m, err := matchModel.GetMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := matchModel.GetPlayersByIdMatch(user.IdLanguage, matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userConfigurationFavoriteSpells, err := userModel.GetConfigurationByIdConfigurationIdUser(user.IdLanguage, user.Id, 1)
	if err != nil {
		h.Slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := matchModel.GetSpellsPlayersByIdMatch(user.IdLanguage, matchId, user.Id, userConfigurationFavoriteSpells.IsEnabled)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.PlayerTable(m, players, spellsPlayers, false).Render(r.Context(), w)
}

func (h *Handler) HandleUsePlayerSpell(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := matchModel.UsePlayerSpellByIdPlayerSpell(idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellPlayer, err := matchModel.GetSpellPlayerByIdSpellsPlayers(user.IdLanguage, idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.Spell(spellPlayer, false).Render(r.Context(), w)
}

func (h *Handler) HandleRemoveRoundRecoveryPlayerSpell(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	h.Slog = h.Slog.With("userId", user.Id)

	vars := mux.Vars(r)
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := matchModel.RemoveRoundRecoverySpellByIdPlayerSpell(idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellPlayer, err := matchModel.GetSpellPlayerByIdSpellsPlayers(user.IdLanguage, idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageMatch.Spell(spellPlayer, false).Render(r.Context(), w)
}
