package tournament

import (
	"context"
)

type Team struct {
	Id        int    `json:"id"`
	IdUser    int    `json:"id_user"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TeamCreate struct {
	IdUser int
	Name   string
}

func GetTeamsByUser(ctx context.Context, idUser int) ([]Team, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, id_user, name, created_at, updated_at FROM tournament_team WHERE id_user = ? ORDER BY name", idUser)
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

func GetTeam(ctx context.Context, id int) (Team, error) {
	t := Team{}
	err := db.QueryRowContext(ctx, "SELECT id, id_user, name, created_at, updated_at FROM tournament_team WHERE id = ?", id).
		Scan(&t.Id, &t.IdUser, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func CreateTeam(ctx context.Context, t TeamCreate) (int, error) {
	response, err := db.ExecContext(ctx, "INSERT INTO tournament_team (id_user, name) VALUES (?, ?)", t.IdUser, t.Name)
	if err != nil {
		return 0, err
	}
	id, _ := response.LastInsertId()
	return int(id), nil
}

func DeleteTeam(ctx context.Context, id int, idUser int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM tournament_team WHERE id = ? AND id_user = ?", id, idUser)
	return err
}
