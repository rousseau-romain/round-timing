package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rousseau-romain/round-timing/model"
	"github.com/rousseau-romain/round-timing/shared/components"
	"github.com/rousseau-romain/round-timing/views/page"

	"github.com/gorilla/mux"
	"github.com/markbates/goth"
)

var PagesNav = []components.NavItem{
	{Name: "Match list", Url: "match"},
}

func getPageNavCustom(user model.User, match model.Match) []components.NavItem {
	var pageNav = PagesNav
	if user.Id != 0 {
		if match.Id != 0 {
			pageNav = append(pageNav, components.NavItem{
				Name: fmt.Sprintf("Match %s (%d)", match.Name, match.Id),
				Url:  fmt.Sprintf("match/%d", match.Id),
			})
		} else {
			lastMatch, err := model.GetLastMatchByUserId(user.Id)
			if err != nil {
				log.Println(err)
				return pageNav
			}
			if lastMatch.Id != 0 {
				pageNav = append(pageNav, components.NavItem{
					Name: fmt.Sprintf("Last match %s (%d)", lastMatch.Name, lastMatch.Id),
					Url:  fmt.Sprintf("match/%d", lastMatch.Id),
				})
			}
		}
	}
	if user.IsAdmin {
		pageNav = append(pageNav, components.NavItem{
			Name: "Admin list users",
			Url:  "admin/user",
		})
	}
	return pageNav
}

func (h *Handler) HandlersNotFound(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	page.NotFoundPage(userOauth2, h.error, PagesNav).Render(r.Context(), w)
}

func (h *Handler) HandlersHome(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	pageNav := PagesNav

	if userOauth2.Name != "" {
		user, err := model.GetUserByOauth2Id(userOauth2.UserID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageNav = getPageNavCustom(user, model.Match{})
		page.HomePage(userOauth2, user, h.error, pageNav).Render(r.Context(), w)
		return
	}
	page.HomePage(goth.User{}, model.User{}, h.error, pageNav).Render(r.Context(), w)
}

func (h *Handler) HandlersProfile(w http.ResponseWriter, r *http.Request) {
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, err := model.GetUserByOauth2Id(userOauth2.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page.ProfilePage(userOauth2, user, h.error, getPageNavCustom(user, model.Match{})).Render(r.Context(), w)
}

func (h *Handler) HandlerStartMatchPage(w http.ResponseWriter, r *http.Request) {
	var idClassGlobal = 13

	userOauth2, _ := h.auth.GetSessionUser(r)
	user, err := model.GetUserByOauth2Id(userOauth2.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	matchId, _ := strconv.Atoi(vars["idMatch"])
	log.Println(matchId)
	match, err := model.GetMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := model.GetPlayersByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(players) == 0 {
		http.Error(w, "Match need players", http.StatusBadRequest)
		return
	}

	if match.Round == 0 {
		classeIds := []int{idClassGlobal}
		idSpellToExclude := []int{134, 135, 136, 137, 139, 140, 141, 142}
		for _, player := range players {
			classeIds = append(classeIds, player.Class.Id)
		}

		spells, err := model.GetSpellsByIdCLass(classeIds, idSpellToExclude)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// maitrise marteau id 138
		if match.MultipleMasteryEnabled == 1 {
			idSpellToExclude = []int{}
		}

		globalSpells, err := model.GetSpellsByIdCLass([]int{idClassGlobal}, idSpellToExclude)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var round = 1
		err = model.UpdateMatch(matchId, model.MatchUpdate{Round: &round})
		if err != nil {
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

	spellsPlayer, err := model.GetSpellsPlayersByIdMatch(matchId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getPageNavCustom(user, model.Match{})

	page.StartMatchPage(userOauth2, user, h.error, getPageNavCustom(user, match), match, players, spellsPlayer).Render(r.Context(), w)
}

func (h *Handler) HandlerResetMatchPage(w http.ResponseWriter, r *http.Request) {
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
	userOauth2, _ := h.auth.GetSessionUser(r)
	user, err := model.GetUserByOauth2Id(userOauth2.UserID)
	vars := mux.Vars(r)
	matchId, _ := strconv.Atoi(vars["idMatch"])
	multipleMasteryEnabled, _ := strconv.Atoi(vars["toggleBool"])

	err = model.UpdateMatch(matchId, model.MatchUpdate{
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

	players, err := model.GetPlayersByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := model.GetSpellsPlayersByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.MatchPageTable(userOauth2, user, h.error, getPageNavCustom(user, match), match, players, spellsPlayers).Render(r.Context(), w)
}

func (h *Handler) HandlerMatchNextRound(w http.ResponseWriter, r *http.Request) {
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

	players, err := model.GetPlayersByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellsPlayers, err := model.GetSpellsPlayersByIdMatch(matchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.PlayerTable(match, players, spellsPlayers).Render(r.Context(), w)
}

func (h *Handler) HandlerUsePlayerSpell(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := model.UsePlayerSpellByIdPlayerSpell(idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellPlayer, err := model.GetSpellPlayerByIdSpellsPlayers(idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.Spell(spellPlayer).Render(r.Context(), w)
}

func (h *Handler) HandlerRemoveRoundRecoveryPlayerSpell(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPlayerSpell, _ := strconv.Atoi(vars["idPlayerSpell"])

	err := model.RemoveRoundRecoverySpellByIdPlayerSpell(idPlayerSpell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spellPlayer, err := model.GetSpellPlayerByIdSpellsPlayers(idPlayerSpell)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.Spell(spellPlayer).Render(r.Context(), w)
}
