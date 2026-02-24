package tournament

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
)

type Tournament struct {
	Id                 int    `json:"id"`
	IdUser             int    `json:"id_user"`
	Name               string `json:"name"`
	NumberPlayerByTeam int    `json:"number_player_by_team"`
	StageType          string `json:"stage_type"`
	Status             string `json:"status"`
	IsArchived         int    `json:"is_archived"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type TournamentCreate struct {
	Name               string
	IdUser             int
	NumberPlayerByTeam int
	StageType          string
}

type TournamentUpdate struct {
	Name      *string
	Status    *string
	StageType *string
}

func GetArchivedTournamentsByIdUser(ctx context.Context, idUser int) ([]Tournament, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "id_user", "name", "number_player_by_team", "stage_type", "status", "is_archived", "created_at", "updated_at").
		From("tournament").
		Where(sb.Equal("id_user", idUser), sb.Equal("is_archived", 1))
	sql, args := sb.Build()

	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []Tournament
	for rows.Next() {
		var t Tournament
		if err := rows.Scan(&t.Id, &t.IdUser, &t.Name, &t.NumberPlayerByTeam, &t.StageType, &t.Status, &t.IsArchived, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return tournaments, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, rows.Err()
}

func GetTournamentsByIdUser(ctx context.Context, idUser int, includeArchived bool) ([]Tournament, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "id_user", "name", "number_player_by_team", "stage_type", "status", "is_archived", "created_at", "updated_at").From("tournament").Where(sb.Equal("id_user", idUser))
	if !includeArchived {
		sb.Where(sb.Equal("is_archived", 0))
	}
	sql, args := sb.Build()

	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []Tournament
	for rows.Next() {
		var t Tournament
		err := rows.Scan(&t.Id, &t.IdUser, &t.Name, &t.NumberPlayerByTeam, &t.StageType, &t.Status, &t.IsArchived, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return tournaments, err
		}
		tournaments = append(tournaments, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tournaments, nil
}

func GetTournament(ctx context.Context, id int) (Tournament, error) {
	t := Tournament{}
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "id_user", "name", "number_player_by_team", "stage_type", "status", "is_archived", "created_at", "updated_at").From("tournament").Where(sb.Equal("id", id))
	sql, args := sb.Build()

	row := db.QueryRowContext(ctx, sql, args...)
	if row.Err() != nil {
		return t, row.Err()
	}
	err := row.Scan(&t.Id, &t.IdUser, &t.Name, &t.NumberPlayerByTeam, &t.StageType, &t.Status, &t.IsArchived, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func GetTournamentOwnerId(ctx context.Context, id int) (int, error) {
	var idUser int
	row := db.QueryRowContext(ctx, "SELECT id_user FROM tournament WHERE id = ?", id)
	if row.Err() != nil {
		return 0, row.Err()
	}
	err := row.Scan(&idUser)
	return idUser, err
}

func CreateTournament(ctx context.Context, t TournamentCreate) (int, error) {
	sb := sqlbuilder.NewInsertBuilder()
	sb.InsertInto("tournament").Cols("id_user", "name", "number_player_by_team", "stage_type").Values(t.IdUser, t.Name, t.NumberPlayerByTeam, t.StageType)
	sql, args := sb.Build()

	response, err := db.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	id, _ := response.LastInsertId()
	return int(id), nil
}

func UpdateTournament(ctx context.Context, id int, t TournamentUpdate) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.Update("tournament").Where(sb.Equal("id", id))
	canUpdate := false

	if t.Name != nil {
		sb.SetMore(sb.Assign("name", *t.Name))
		canUpdate = true
	}
	if t.Status != nil {
		sb.SetMore(sb.Assign("status", *t.Status))
		canUpdate = true
	}
	if t.StageType != nil {
		sb.SetMore(sb.Assign("stage_type", *t.StageType))
		canUpdate = true
	}
	if !canUpdate {
		return nil
	}
	sql, args := sb.Build()
	_, err := db.ExecContext(ctx, sql, args...)
	return err
}

func ArchiveTournament(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "UPDATE tournament SET is_archived = 1 WHERE id = ?", id)
	return err
}

func UnarchiveTournament(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "UPDATE tournament SET is_archived = 0 WHERE id = ?", id)
	return err
}

func DeleteTournament(ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "DELETE FROM tournament WHERE id = ?", id)
	return err
}
