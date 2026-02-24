package tournament

import (
	"context"
	"database/sql"
	"fmt"
)

type TournamentMatch struct {
	Id           int           `json:"id"`
	IdTournament int           `json:"id_tournament"`
	IdTeam1      sql.NullInt64 `json:"id_team1"`
	IdTeam2      sql.NullInt64 `json:"id_team2"`
	IdTeamWinner sql.NullInt64 `json:"id_team_winner"`
	BoFormat     int           `json:"bo_format"`
	ScoreTeam1   int           `json:"score_team1"`
	ScoreTeam2   int           `json:"score_team2"`
	KillsTeam1   int           `json:"kills_team1"`
	KillsTeam2   int           `json:"kills_team2"`
	Round        int           `json:"round"`
	Position     int           `json:"position"`
	Status       string        `json:"status"`
	CreatedAt    string        `json:"created_at"`
}

type TournamentMatchCreate struct {
	IdTournament int
	IdTeam1      int
	IdTeam2      int
	BoFormat     int
	Round        int
	Position     int
}

type TournamentMatchWithNames struct {
	TournamentMatch
	Team1Name  string
	Team2Name  string
	WinnerName string
}

func GetMatchesByTournament(ctx context.Context, idTournament int) ([]TournamentMatchWithNames, error) {
	query := `
		SELECT
			m.id, m.id_tournament, m.id_team1, m.id_team2,
			m.id_team_winner, m.bo_format, m.score_team1, m.score_team2,
			m.kills_team1, m.kills_team2,
			m.round, m.position, m.status, m.created_at,
			COALESCE(t1.name, ''), COALESCE(t2.name, ''), COALESCE(tw.name, '')
		FROM tournament_match m
		LEFT JOIN tournament_team t1 ON t1.id = m.id_team1
		LEFT JOIN tournament_team t2 ON t2.id = m.id_team2
		LEFT JOIN tournament_team tw ON tw.id = m.id_team_winner
		WHERE m.id_tournament = ?
		ORDER BY m.round, m.position
	`
	rows, err := db.QueryContext(ctx, query, idTournament)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []TournamentMatchWithNames
	for rows.Next() {
		var m TournamentMatchWithNames
		if err := rows.Scan(&m.Id, &m.IdTournament, &m.IdTeam1, &m.IdTeam2,
			&m.IdTeamWinner, &m.BoFormat, &m.ScoreTeam1, &m.ScoreTeam2,
			&m.KillsTeam1, &m.KillsTeam2,
			&m.Round, &m.Position, &m.Status, &m.CreatedAt,
			&m.Team1Name, &m.Team2Name, &m.WinnerName); err != nil {
			return matches, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func GetMatchesByTournamentAndRound(ctx context.Context, idTournament int, round int) ([]TournamentMatchWithNames, error) {
	query := `
		SELECT
			m.id, m.id_tournament, m.id_team1, m.id_team2,
			m.id_team_winner, m.bo_format, m.score_team1, m.score_team2,
			m.kills_team1, m.kills_team2,
			m.round, m.position, m.status, m.created_at,
			COALESCE(t1.name, ''), COALESCE(t2.name, ''), COALESCE(tw.name, '')
		FROM tournament_match m
		LEFT JOIN tournament_team t1 ON t1.id = m.id_team1
		LEFT JOIN tournament_team t2 ON t2.id = m.id_team2
		LEFT JOIN tournament_team tw ON tw.id = m.id_team_winner
		WHERE m.id_tournament = ? AND m.round = ?
		ORDER BY m.position
	`
	rows, err := db.QueryContext(ctx, query, idTournament, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []TournamentMatchWithNames
	for rows.Next() {
		var m TournamentMatchWithNames
		if err := rows.Scan(&m.Id, &m.IdTournament, &m.IdTeam1, &m.IdTeam2,
			&m.IdTeamWinner, &m.BoFormat, &m.ScoreTeam1, &m.ScoreTeam2,
			&m.KillsTeam1, &m.KillsTeam2,
			&m.Round, &m.Position, &m.Status, &m.CreatedAt,
			&m.Team1Name, &m.Team2Name, &m.WinnerName); err != nil {
			return matches, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

const matchColumns = "id, id_tournament, id_team1, id_team2, id_team_winner, bo_format, score_team1, score_team2, kills_team1, kills_team2, round, position, status, created_at"

func scanTournamentMatch(scanner interface{ Scan(...any) error }) (TournamentMatch, error) {
	var m TournamentMatch
	err := scanner.Scan(&m.Id, &m.IdTournament, &m.IdTeam1, &m.IdTeam2,
		&m.IdTeamWinner, &m.BoFormat, &m.ScoreTeam1, &m.ScoreTeam2,
		&m.KillsTeam1, &m.KillsTeam2,
		&m.Round, &m.Position, &m.Status, &m.CreatedAt)
	return m, err
}

func GetTournamentMatch(ctx context.Context, id int) (TournamentMatch, error) {
	return scanTournamentMatch(db.QueryRowContext(ctx,
		"SELECT "+matchColumns+" FROM tournament_match WHERE id = ?", id))
}

func CreateTournamentMatch(ctx context.Context, m TournamentMatchCreate) (int, error) {
	response, err := db.ExecContext(ctx,
		"INSERT INTO tournament_match (id_tournament, id_team1, id_team2, bo_format, round, position) VALUES (?, ?, ?, ?, ?, ?)",
		m.IdTournament, m.IdTeam1, m.IdTeam2, m.BoFormat, m.Round, m.Position)
	if err != nil {
		return 0, err
	}
	id, _ := response.LastInsertId()
	return int(id), nil
}

func UpdateMatchBo(ctx context.Context, idMatch int, boFormat int) error {
	_, err := db.ExecContext(ctx, "UPDATE tournament_match SET bo_format = ? WHERE id = ?", boFormat, idMatch)
	return err
}

func UpdateMatchsBo(ctx context.Context, idTournament int, round int, boFormat int) error {
	_, err := db.ExecContext(ctx, "UPDATE tournament_match SET bo_format = ? WHERE id_tournament = ? AND round = ?", boFormat, idTournament, round)
	return err
}

func UpdateMatchScore(ctx context.Context, idMatch int, scoreTeam1 int, scoreTeam2 int) error {
	_, err := db.ExecContext(ctx, "UPDATE tournament_match SET score_team1 = ?, score_team2 = ? WHERE id = ?", scoreTeam1, scoreTeam2, idMatch)
	return err
}

func UpdateMatchWinner(ctx context.Context, idMatch int, idTeamWinner *int, status string) error {
	if idTeamWinner != nil {
		_, err := db.ExecContext(ctx, "UPDATE tournament_match SET id_team_winner = ?, status = ? WHERE id = ?", *idTeamWinner, status, idMatch)
		return err
	}
	_, err := db.ExecContext(ctx, "UPDATE tournament_match SET id_team_winner = NULL, status = ? WHERE id = ?", status, idMatch)
	return err
}

func GetTeamWinningRound(ctx context.Context, idTournament int, idTeam int) (int, error) {
	var round int
	err := db.QueryRowContext(ctx,
		"SELECT COALESCE(MAX(round), 0) FROM tournament_match WHERE id_tournament = ? AND id_team_winner = ? AND status = 'finished'",
		idTournament, idTeam).Scan(&round)
	return round, err
}

func GetMatchCountInRound(ctx context.Context, idTournament int, round int) (int, error) {
	var count int
	err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM tournament_match WHERE id_tournament = ? AND round = ?",
		idTournament, round).Scan(&count)
	return count, err
}

func GetAvailableTeamsForTournament(ctx context.Context, idTournament int, idUser int) ([]Team, error) {
	// Check if any finished matches exist
	var finishedCount int
	err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM tournament_match WHERE id_tournament = ? AND status = 'finished'",
		idTournament).Scan(&finishedCount)
	if err != nil {
		return nil, err
	}

	// No finished matches yet: return all user teams not already in a match for this tournament
	if finishedCount == 0 {
		query := `
			SELECT t.id, t.id_user, t.name, t.created_at, t.updated_at
			FROM tournament_team t
			WHERE t.id_user = ?
			AND t.id NOT IN (
				SELECT m.id_team1 FROM tournament_match m WHERE m.id_tournament = ? AND m.id_team1 IS NOT NULL
				UNION
				SELECT m.id_team2 FROM tournament_match m WHERE m.id_tournament = ? AND m.id_team2 IS NOT NULL
			)
			ORDER BY t.name
		`
		rows, err := db.QueryContext(ctx, query, idUser, idTournament, idTournament)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var teams []Team
		for rows.Next() {
			var t Team
			if err := rows.Scan(&t.Id, &t.IdUser, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
				return teams, err
			}
			teams = append(teams, t)
		}
		return teams, rows.Err()
	}

	// Return winners that are not already placed in a higher-round match
	query := `
		SELECT DISTINCT t.id, t.id_user, t.name, t.created_at, t.updated_at
		FROM tournament_match m
		JOIN tournament_team t ON t.id = m.id_team_winner
		WHERE m.id_tournament = ? AND m.status = 'finished'
		AND t.id NOT IN (
			SELECT m2.id_team1 FROM tournament_match m2
			WHERE m2.id_tournament = ? AND m2.id_team1 IS NOT NULL
			AND m2.round > m.round
			UNION
			SELECT m2.id_team2 FROM tournament_match m2
			WHERE m2.id_tournament = ? AND m2.id_team2 IS NOT NULL
			AND m2.round > m.round
		)
		ORDER BY t.name
	`
	rows, err := db.QueryContext(ctx, query, idTournament, idTournament, idTournament)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.Id, &t.IdUser, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return teams, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func GetMatchByRoundAndPosition(ctx context.Context, idTournament int, round int, position int) (TournamentMatch, error) {
	return scanTournamentMatch(db.QueryRowContext(ctx,
		"SELECT "+matchColumns+" FROM tournament_match WHERE id_tournament = ? AND round = ? AND position = ?",
		idTournament, round, position))
}

func GetMaxRound(ctx context.Context, idTournament int) (int, error) {
	var round int
	err := db.QueryRowContext(ctx,
		"SELECT COALESCE(MAX(round), 0) FROM tournament_match WHERE id_tournament = ?",
		idTournament).Scan(&round)
	return round, err
}

func IsAllMatchesFinishedInRound(ctx context.Context, idTournament int, round int) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM tournament_match WHERE id_tournament = ? AND round = ? AND status != 'finished'",
		idTournament, round).Scan(&count)
	return count == 0, err
}

func GetWinnersOfRound(ctx context.Context, idTournament int, round int) ([]int, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id_team_winner FROM tournament_match
		WHERE id_tournament = ? AND round = ? AND status = 'finished' AND id_team_winner IS NOT NULL`,
		idTournament, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var winners []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return winners, err
		}
		winners = append(winners, id)
	}
	return winners, rows.Err()
}

func IsTeamInRound(ctx context.Context, idTournament int, idTeam int, round int) (bool, error) {
	var count int
	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tournament_match
		WHERE id_tournament = ? AND round = ?
		AND (id_team1 = ? OR id_team2 = ?)`,
		idTournament, round, idTeam, idTeam).Scan(&count)
	return count > 0, err
}

func DeleteTournamentMatch(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM tournament_match WHERE id = ?", id)
	return err
}

func DeleteTournamentMatchesByTournament(ctx context.Context, idTournament int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM tournament_match WHERE id_tournament = ?", idTournament)
	return err
}

func incrementMatchField(ctx context.Context, idMatch int, fieldPrefix string, team int, delta int) error {
	col := fieldPrefix + "_team1"
	if team == 2 {
		col = fieldPrefix + "_team2"
	}
	query := fmt.Sprintf(
		"UPDATE tournament_match SET %s = GREATEST(0, %s + ?) WHERE id = ?", col, col)
	_, err := db.ExecContext(ctx, query, delta, idMatch)
	return err
}

func IncrementMatchKills(ctx context.Context, idMatch int, team int, delta int) error {
	return incrementMatchField(ctx, idMatch, "kills", team, delta)
}

func IncrementMatchScore(ctx context.Context, idMatch int, team int, delta int) error {
	return incrementMatchField(ctx, idMatch, "score", team, delta)
}

func UpdateMatchStatus(ctx context.Context, idMatch int, status string) error {
	_, err := db.ExecContext(ctx, "UPDATE tournament_match SET status = ? WHERE id = ?", status, idMatch)
	return err
}
