package tournament

import (
	"context"
)

type TeamPlayer struct {
	Id        int    `json:"id"`
	IdUser    int    `json:"id_user"`
	IdTeam    int    `json:"id_team"`
	IdPlayer  int    `json:"id_player"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TeamPlayerWithNames struct {
	TeamPlayer
	TeamName   string
	PlayerName string
}

type TeamPlayerCreate struct {
	IdUser   int
	IdTeam   int
	IdPlayer int
}

func GetTeamPlayersByTeam(ctx context.Context, idTeam int) ([]TeamPlayerWithNames, error) {
	query := `
		SELECT tp.id, tp.id_user, tp.id_team, tp.id_player, tp.created_at, tp.updated_at, t.name, p.name
		FROM tournament_team_player tp
		JOIN tournament_team t ON t.id = tp.id_team
		JOIN tournament_player p ON p.id = tp.id_player
		WHERE tp.id_team = ?
		ORDER BY p.name
	`
	rows, err := db.QueryContext(ctx, query, idTeam)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamPlayers []TeamPlayerWithNames
	for rows.Next() {
		var tp TeamPlayerWithNames
		if err := rows.Scan(&tp.Id, &tp.IdUser, &tp.IdTeam, &tp.IdPlayer, &tp.CreatedAt, &tp.UpdatedAt, &tp.TeamName, &tp.PlayerName); err != nil {
			return teamPlayers, err
		}
		teamPlayers = append(teamPlayers, tp)
	}
	return teamPlayers, rows.Err()
}

func GetTeamPlayersByUser(ctx context.Context, idUser int) ([]TeamPlayerWithNames, error) {
	query := `
		SELECT tp.id, tp.id_user, tp.id_team, tp.id_player, tp.created_at, tp.updated_at, t.name, p.name
		FROM tournament_team_player tp
		JOIN tournament_team t ON t.id = tp.id_team
		JOIN tournament_player p ON p.id = tp.id_player
		WHERE tp.id_user = ?
		ORDER BY t.name, p.name
	`
	rows, err := db.QueryContext(ctx, query, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamPlayers []TeamPlayerWithNames
	for rows.Next() {
		var tp TeamPlayerWithNames
		if err := rows.Scan(&tp.Id, &tp.IdUser, &tp.IdTeam, &tp.IdPlayer, &tp.CreatedAt, &tp.UpdatedAt, &tp.TeamName, &tp.PlayerName); err != nil {
			return teamPlayers, err
		}
		teamPlayers = append(teamPlayers, tp)
	}
	return teamPlayers, rows.Err()
}

func CreateTeamPlayer(ctx context.Context, tp TeamPlayerCreate) (int, error) {
	response, err := db.ExecContext(ctx, "INSERT INTO tournament_team_player (id_user, id_team, id_player) VALUES (?, ?, ?)", tp.IdUser, tp.IdTeam, tp.IdPlayer)
	if err != nil {
		return 0, err
	}
	id, _ := response.LastInsertId()
	return int(id), nil
}

func DeleteTeamPlayer(ctx context.Context, id int, idUser int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM tournament_team_player WHERE id = ? AND id_user = ?", id, idUser)
	return err
}
