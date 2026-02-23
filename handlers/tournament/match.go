package tournament

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rousseau-romain/round-timing/handlers"
	tournamentModel "github.com/rousseau-romain/round-timing/model/tournament"
	httpError "github.com/rousseau-romain/round-timing/pkg/httperror"
	"github.com/rousseau-romain/round-timing/service/auth"
	pageTournament "github.com/rousseau-romain/round-timing/views/page/tournament"
)

func stageTypeFromWinnerCount(count int) string {
	switch {
	case count >= 32:
		return "round_32"
	case count >= 16:
		return "round_16"
	case count >= 8:
		return "quarter_final"
	case count >= 4:
		return "semi_final"
	default:
		return "final"
	}
}

func (h *Handler) HandleCreateMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	idTeam1, err := strconv.Atoi(r.FormValue("id_team1"))
	if err != nil {
		handlers.RenderComponentError("Invalid team 1", []string{"Invalid team 1"}, http.StatusBadRequest, w, r)
		return
	}
	idTeam2, err := strconv.Atoi(r.FormValue("id_team2"))
	if err != nil {
		handlers.RenderComponentError("Invalid team 2", []string{"Invalid team 2"}, http.StatusBadRequest, w, r)
		return
	}
	if idTeam1 == idTeam2 {
		handlers.RenderComponentError("Teams must be different", []string{"Teams must be different"}, http.StatusBadRequest, w, r)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	team1Valid, team2Valid := false, false
	for _, t := range availableTeams {
		if t.Id == idTeam1 {
			team1Valid = true
		}
		if t.Id == idTeam2 {
			team2Valid = true
		}
	}
	if !team1Valid || !team2Valid {
		handlers.RenderComponentError("Team is not an available winner", []string{"Team is not an available winner"}, http.StatusBadRequest, w, r)
		return
	}

	boFormat, err := strconv.Atoi(r.FormValue("bo_format"))
	if err != nil || boFormat < 1 {
		boFormat = 1
	}

	// Auto-compute round from the teams' winning rounds
	round1, err := tournamentModel.GetTeamWinningRound(r.Context(), idTournament, idTeam1)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	round2, err := tournamentModel.GetTeamWinningRound(r.Context(), idTournament, idTeam2)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	round := max(round1, round2) + 1

	// Check that neither team is already in a match for this round
	team1InRound, err := tournamentModel.IsTeamInRound(r.Context(), idTournament, idTeam1, round)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	team2InRound, err := tournamentModel.IsTeamInRound(r.Context(), idTournament, idTeam2, round)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if team1InRound || team2InRound {
		handlers.RenderComponentError("Team is already in a match for this round", []string{"Team is already in a match for this round"}, http.StatusBadRequest, w, r)
		return
	}

	// Auto-compute position within the round
	position, err := tournamentModel.GetMatchCountInRound(r.Context(), idTournament, round)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	matchId, err := tournamentModel.CreateTournamentMatch(r.Context(), tournamentModel.TournamentMatchCreate{
		IdTournament: idTournament,
		IdTeam1:      idTeam1,
		IdTeam2:      idTeam2,
		BoFormat:     boFormat,
		Round:        round,
		Position:     position,
	})
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("match created", "tournamentId", idTournament, "matchId", matchId)

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	// Re-fetch available teams after match creation
	availableTeams, err = tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	w.Header().Set("HX-Trigger", "tournamentStatusChanged")
	pageTournament.MatchesSection(t, matches, availableTeams).Render(r.Context(), w)
}

func (h *Handler) HandleUpdateMatchScore(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	scoreTeam1, err := strconv.Atoi(r.FormValue("score_team1"))
	if err != nil {
		http.Error(w, "Invalid score team 1", http.StatusBadRequest)
		return
	}
	scoreTeam2, err := strconv.Atoi(r.FormValue("score_team2"))
	if err != nil {
		http.Error(w, "Invalid score team 2", http.StatusBadRequest)
		return
	}

	if err := tournamentModel.UpdateMatchScore(r.Context(), idMatch, scoreTeam1, scoreTeam2); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	m, err := tournamentModel.GetTournamentMatch(r.Context(), idMatch)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	var winnerId *int
	winsNeeded := (m.BoFormat / 2) + 1
	if scoreTeam1 >= winsNeeded && m.IdTeam1.Valid {
		id := int(m.IdTeam1.Int64)
		winnerId = &id
		if err := tournamentModel.UpdateMatchWinner(r.Context(), idMatch, winnerId, "finished"); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
	} else if scoreTeam2 >= winsNeeded && m.IdTeam2.Valid {
		id := int(m.IdTeam2.Int64)
		winnerId = &id
		if err := tournamentModel.UpdateMatchWinner(r.Context(), idMatch, winnerId, "finished"); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
	} else if scoreTeam1 > 0 || scoreTeam2 > 0 {
		if err := tournamentModel.UpdateMatchWinner(r.Context(), idMatch, nil, "in_progress"); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
	}

	logger.Info("match score updated", "matchId", idMatch, "scoreTeam1", scoreTeam1, "scoreTeam2", scoreTeam2)

	// Update tournament status based on match states
	currentTournament, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if currentTournament.Status == "draft" && (scoreTeam1 > 0 || scoreTeam2 > 0) {
		status := "in_progress"
		if err := tournamentModel.UpdateTournament(r.Context(), idTournament, tournamentModel.TournamentUpdate{Status: &status}); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
		currentTournament.Status = "in_progress"
	}
	matchCount, err := tournamentModel.GetMatchCountInRound(r.Context(), idTournament, m.Round)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if matchCount == 1 {
		var newStatus string
		if winnerId != nil {
			newStatus = "finished"
		} else if currentTournament.Status == "finished" {
			newStatus = "in_progress"
		}
		if newStatus != "" && newStatus != currentTournament.Status {
			if err := tournamentModel.UpdateTournament(r.Context(), idTournament, tournamentModel.TournamentUpdate{Status: &newStatus}); err != nil {
				logger.Error(err.Error())
				httpError.InternalError(w)
				return
			}
		}
	}

	// Auto-create final match if exactly 2 winners remain unmatched
	if winnerId != nil {
		allRoundFinished, err := tournamentModel.IsAllMatchesFinishedInRound(r.Context(), idTournament, m.Round)
		if err == nil && allRoundFinished {
			nextAvailable, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
			if err == nil && len(nextAvailable) == 2 {
				t1Round, _ := tournamentModel.GetTeamWinningRound(r.Context(), idTournament, nextAvailable[0].Id)
				t2Round, _ := tournamentModel.GetTeamWinningRound(r.Context(), idTournament, nextAvailable[1].Id)
				finalRound := max(t1Round, t2Round) + 1
				if count, _ := tournamentModel.GetMatchCountInRound(r.Context(), idTournament, finalRound); count == 0 {
					if _, err := tournamentModel.CreateTournamentMatch(r.Context(), tournamentModel.TournamentMatchCreate{
						IdTournament: idTournament,
						IdTeam1:      nextAvailable[0].Id,
						IdTeam2:      nextAvailable[1].Id,
						BoFormat:     1,
						Round:        finalRound,
						Position:     0,
					}); err != nil {
						logger.Error("auto-final: failed to create match", "error", err.Error())
					} else {
						logger.Info("auto-created final match", "tournamentId", idTournament, "round", finalRound)
						finalStageType := "final"
						_ = tournamentModel.UpdateTournament(r.Context(), idTournament, tournamentModel.TournamentUpdate{StageType: &finalStageType})
					}
				}
			}
		}
	}

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	w.Header().Set("HX-Trigger", "tournamentStatusChanged")
	pageTournament.MatchesSection(t, matches, availableTeams).Render(r.Context(), w)
}

func (h *Handler) HandleEditMatchBo(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	boFormat, err := strconv.Atoi(r.FormValue("bo_format"))
	if err != nil || boFormat < 1 {
		http.Error(w, "Invalid BO format", http.StatusBadRequest)
		return
	}

	if err := tournamentModel.UpdateMatchBo(r.Context(), idMatch, boFormat); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("match BO updated", "matchId", idMatch, "boFormat", boFormat)

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	w.Header().Set("HX-Trigger", "tournamentStatusChanged")
	pageTournament.MatchesSection(t, matches, availableTeams).Render(r.Context(), w)
}

func (h *Handler) generateNextRound(w http.ResponseWriter, r *http.Request, shuffle bool) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	currentRound, err := tournamentModel.GetMaxRound(r.Context(), idTournament)
	if err != nil || currentRound == 0 {
		handlers.RenderComponentError("No matches found", []string{"No matches found"}, http.StatusBadRequest, w, r)
		return
	}

	allFinished, err := tournamentModel.IsAllMatchesFinishedInRound(r.Context(), idTournament, currentRound)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if !allFinished {
		handlers.RenderComponentError("All matches in the current round must be finished", []string{"All matches in the current round must be finished"}, http.StatusBadRequest, w, r)
		return
	}

	winners, err := tournamentModel.GetWinnersOfRound(r.Context(), idTournament, currentRound)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if len(winners) < 2 {
		handlers.RenderComponentError("Not enough winners to generate next round", []string{"Not enough winners to generate next round"}, http.StatusBadRequest, w, r)
		return
	}

	if shuffle {
		rand.Shuffle(len(winners), func(i, j int) {
			winners[i], winners[j] = winners[j], winners[i]
		})
	}

	nextRound := currentRound + 1
	for i := 0; i+1 < len(winners); i += 2 {
		_, err := tournamentModel.CreateTournamentMatch(r.Context(), tournamentModel.TournamentMatchCreate{
			IdTournament: idTournament,
			IdTeam1:      winners[i],
			IdTeam2:      winners[i+1],
			BoFormat:     1,
			Round:        nextRound,
			Position:     i / 2,
		})
		if err != nil {
			logger.Error("failed to create next round match", "error", err.Error())
			httpError.InternalError(w)
			return
		}
	}

	logger.Info("generated next round", "tournamentId", idTournament, "round", nextRound, "matches", len(winners)/2, "shuffle", shuffle)

	nextStageType := stageTypeFromWinnerCount(len(winners))
	if err := tournamentModel.UpdateTournament(r.Context(), idTournament, tournamentModel.TournamentUpdate{StageType: &nextStageType}); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	w.Header().Set("HX-Trigger", "tournamentStatusChanged")
	pageTournament.MatchesSection(t, matches, availableTeams).Render(r.Context(), w)
}

func (h *Handler) HandleGenerateNextRound(w http.ResponseWriter, r *http.Request) {
	h.generateNextRound(w, r, true)
}

func (h *Handler) HandleGenerateNextRoundOrdered(w http.ResponseWriter, r *http.Request) {
	h.generateNextRound(w, r, false)
}

func (h *Handler) HandleDeleteMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if t.IsArchived == 1 {
		handlers.RenderComponentError("Cannot delete match from an archived tournament", []string{"Cannot delete match from an archived tournament"}, http.StatusBadRequest, w, r)
		return
	}

	currentRound, err := tournamentModel.GetMaxRound(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	m, err := tournamentModel.GetTournamentMatch(r.Context(), idMatch)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if m.Round != currentRound {
		handlers.RenderComponentError("Cannot delete match from a previous round", []string{"Cannot delete match from a previous round"}, http.StatusBadRequest, w, r)
		return
	}

	if err := tournamentModel.DeleteTournamentMatch(r.Context(), idMatch); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	update := tournamentModel.TournamentUpdate{}
	if t.Status == "finished" {
		status := "in_progress"
		update.Status = &status
	}

	// Recalculate StageType if the deleted match was the last in its round
	remainingInRound, err := tournamentModel.GetMatchCountInRound(r.Context(), idTournament, m.Round)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if remainingInRound == 0 {
		newMaxRound, err := tournamentModel.GetMaxRound(r.Context(), idTournament)
		if err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
		if newMaxRound > 0 {
			prevRoundCount, err := tournamentModel.GetMatchCountInRound(r.Context(), idTournament, newMaxRound)
			if err != nil {
				logger.Error(err.Error())
				httpError.InternalError(w)
				return
			}
			newStageType := stageTypeFromWinnerCount(prevRoundCount * 2)
			if newStageType != t.StageType {
				update.StageType = &newStageType
			}
		}
	}

	if err := tournamentModel.UpdateTournament(r.Context(), idTournament, update); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("match deleted", "matchId", idMatch)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleUpdateMatchKills(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	team, err := strconv.Atoi(r.FormValue("team"))
	if err != nil || (team != 1 && team != 2) {
		http.Error(w, "Invalid team", http.StatusBadRequest)
		return
	}
	delta, err := strconv.Atoi(r.FormValue("delta"))
	if err != nil || (delta != 1 && delta != -1) {
		http.Error(w, "Invalid delta", http.StatusBadRequest)
		return
	}

	if err := tournamentModel.IncrementMatchKills(r.Context(), idMatch, team, delta); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	logger.Info("match kills updated", "matchId", idMatch, "team", team, "delta", delta)

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	pageTournament.MatchesSection(t, matches, availableTeams).Render(r.Context(), w)
}

func (h *Handler) HandleIncrementMatchScore(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idMatch, _ := strconv.Atoi(vars["idMatch"])
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		handlers.RespondWithError(w, r, h.Slog, err, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	team, err := strconv.Atoi(r.FormValue("team"))
	if err != nil || (team != 1 && team != 2) {
		http.Error(w, "Invalid team", http.StatusBadRequest)
		return
	}
	delta, err := strconv.Atoi(r.FormValue("delta"))
	if err != nil || (delta != 1 && delta != -1) {
		http.Error(w, "Invalid delta", http.StatusBadRequest)
		return
	}

	if err := tournamentModel.IncrementMatchScore(r.Context(), idMatch, team, delta); err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	m, err := tournamentModel.GetTournamentMatch(r.Context(), idMatch)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	scoreTeam1 := m.ScoreTeam1
	scoreTeam2 := m.ScoreTeam2

	var winnerId *int
	winsNeeded := (m.BoFormat / 2) + 1
	if scoreTeam1 >= winsNeeded && m.IdTeam1.Valid {
		id := int(m.IdTeam1.Int64)
		winnerId = &id
		if err := tournamentModel.UpdateMatchWinner(r.Context(), idMatch, winnerId, "finished"); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
	} else if scoreTeam2 >= winsNeeded && m.IdTeam2.Valid {
		id := int(m.IdTeam2.Int64)
		winnerId = &id
		if err := tournamentModel.UpdateMatchWinner(r.Context(), idMatch, winnerId, "finished"); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
	} else if scoreTeam1 > 0 || scoreTeam2 > 0 {
		if err := tournamentModel.UpdateMatchWinner(r.Context(), idMatch, nil, "in_progress"); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
	} else {
		if err := tournamentModel.UpdateMatchWinner(r.Context(), idMatch, nil, "pending"); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
	}

	logger.Info("match score incremented", "matchId", idMatch, "team", team, "delta", delta)

	currentTournament, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if currentTournament.Status == "draft" && (scoreTeam1 > 0 || scoreTeam2 > 0) {
		status := "in_progress"
		if err := tournamentModel.UpdateTournament(r.Context(), idTournament, tournamentModel.TournamentUpdate{Status: &status}); err != nil {
			logger.Error(err.Error())
			httpError.InternalError(w)
			return
		}
		currentTournament.Status = "in_progress"
	}
	matchCount, err := tournamentModel.GetMatchCountInRound(r.Context(), idTournament, m.Round)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}
	if matchCount == 1 {
		var newStatus string
		if winnerId != nil {
			newStatus = "finished"
		} else if currentTournament.Status == "finished" {
			newStatus = "in_progress"
		}
		if newStatus != "" && newStatus != currentTournament.Status {
			if err := tournamentModel.UpdateTournament(r.Context(), idTournament, tournamentModel.TournamentUpdate{Status: &newStatus}); err != nil {
				logger.Error(err.Error())
				httpError.InternalError(w)
				return
			}
		}
	}

	if winnerId != nil {
		allRoundFinished, err := tournamentModel.IsAllMatchesFinishedInRound(r.Context(), idTournament, m.Round)
		if err == nil && allRoundFinished {
			nextAvailable, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
			if err == nil && len(nextAvailable) == 2 {
				t1Round, _ := tournamentModel.GetTeamWinningRound(r.Context(), idTournament, nextAvailable[0].Id)
				t2Round, _ := tournamentModel.GetTeamWinningRound(r.Context(), idTournament, nextAvailable[1].Id)
				finalRound := max(t1Round, t2Round) + 1
				if count, _ := tournamentModel.GetMatchCountInRound(r.Context(), idTournament, finalRound); count == 0 {
					if _, err := tournamentModel.CreateTournamentMatch(r.Context(), tournamentModel.TournamentMatchCreate{
						IdTournament: idTournament,
						IdTeam1:      nextAvailable[0].Id,
						IdTeam2:      nextAvailable[1].Id,
						BoFormat:     1,
						Round:        finalRound,
						Position:     0,
					}); err != nil {
						logger.Error("auto-final: failed to create match", "error", err.Error())
					} else {
						logger.Info("auto-created final match", "tournamentId", idTournament, "round", finalRound)
						finalStageType := "final"
						_ = tournamentModel.UpdateTournament(r.Context(), idTournament, tournamentModel.TournamentUpdate{StageType: &finalStageType})
					}
				}
			}
		}
	}

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	w.Header().Set("HX-Trigger", "tournamentStatusChanged")
	pageTournament.MatchesSection(t, matches, availableTeams).Render(r.Context(), w)
}

func (h *Handler) HandleMatchesSection(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromRequest(r)
	logger := h.Slog.With("userId", user.Id)
	vars := mux.Vars(r)
	idTournament, _ := strconv.Atoi(vars["idTournament"])

	t, err := tournamentModel.GetTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	matches, err := tournamentModel.GetMatchesByTournament(r.Context(), idTournament)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	availableTeams, err := tournamentModel.GetAvailableTeamsForTournament(r.Context(), idTournament, user.Id)
	if err != nil {
		logger.Error(err.Error())
		httpError.InternalError(w)
		return
	}

	w.Header().Set("HX-Trigger", "tournamentStatusChanged")
	pageTournament.MatchesSection(t, matches, availableTeams).Render(r.Context(), w)
}
