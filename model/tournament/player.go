package tournament

import (
	"context"
)

type Player struct {
	Id        int    `json:"id"`
	IdUser    int    `json:"id_user"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PlayerCreate struct {
	IdUser int
	Name   string
}

func GetPlayersByUser(ctx context.Context, idUser int) ([]Player, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, id_user, name, created_at, updated_at FROM tournament_player WHERE id_user = ? ORDER BY name", idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Id, &p.IdUser, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return players, err
		}
		players = append(players, p)
	}
	return players, rows.Err()
}

func GetPlayer(ctx context.Context, id int) (Player, error) {
	p := Player{}
	err := db.QueryRowContext(ctx, "SELECT id, id_user, name, created_at, updated_at FROM tournament_player WHERE id = ?", id).
		Scan(&p.Id, &p.IdUser, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func CreatePlayer(ctx context.Context, p PlayerCreate) (int, error) {
	response, err := db.ExecContext(ctx, "INSERT INTO tournament_player (id_user, name) VALUES (?, ?)", p.IdUser, p.Name)
	if err != nil {
		return 0, err
	}
	id, _ := response.LastInsertId()
	return int(id), nil
}

func DeletePlayer(ctx context.Context, id int, idUser int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM tournament_player WHERE id = ? AND id_user = ?", id, idUser)
	return err
}
