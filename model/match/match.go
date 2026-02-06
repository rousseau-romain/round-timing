package match

import (
	"context"
	"errors"

	"github.com/huandu/go-sqlbuilder"
)

type Match struct {
	Id                     int    `json:"id"`
	IdUser                 int    `json:"id_user"`
	Name                   string `json:"name"`
	Round                  int    `json:"round"`
	MultipleMasteryEnabled int    `json:"multiple_mastery_enabled"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}

type MatchCreate struct {
	Name   string
	IdUser int
}

type MatchUpdate struct {
	Name                   *string
	IdUser                 *int
	Round                  *int
	MultipleMasteryEnabled *int
}

func GetMatchsByIdUser(ctx context.Context, idUser int) ([]Match, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "id_user", "name", "round", "multiple_mastery_enabled", "created_at", "updated_at").From("`match`").Where(sb.Equal("id_user", idUser))
	sql, args := sb.Build()

	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchs []Match

	for rows.Next() {
		var match Match
		err := rows.Scan(&match.Id, &match.IdUser, &match.Name, &match.Round, &match.MultipleMasteryEnabled, &match.CreatedAt, &match.UpdatedAt)
		if err != nil {
			return matchs, err
		}
		matchs = append(matchs, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matchs, nil
}

func GetMatch(ctx context.Context, idMatch int) (Match, error) {
	match := Match{}

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "id_user", "name", "round", "multiple_mastery_enabled", "created_at", "updated_at").From("`match`").Where(sb.Equal("id", idMatch))
	sql, args := sb.Build()

	row := db.QueryRowContext(ctx, sql, args...)
	if row.Err() != nil {
		return match, row.Err()
	}
	err := row.Scan(&match.Id, &match.IdUser, &match.Name, &match.Round, &match.MultipleMasteryEnabled, &match.CreatedAt, &match.UpdatedAt)

	return match, err
}

func GetLastMatchByUserId(ctx context.Context, idUser int) (Match, error) {
	match := Match{}

	sql := `
		SELECT
			id,
			id_user,
			name,
			round,
			multiple_mastery_enabled,
			created_at,
			updated_at
		FROM ` + "`match`" + `
		WHERE id_user = ?
		AND updated_at = (SELECT MAX(updated_at) FROM ` + "`match`" + ` WHERE id_user = ?)
	`
	row := db.QueryRowContext(ctx, sql, idUser, idUser)
	if row.Err() != nil {
		return match, row.Err()
	}
	err := row.Scan(&match.Id, &match.IdUser, &match.Name, &match.Round, &match.MultipleMasteryEnabled, &match.CreatedAt, &match.UpdatedAt)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return match, nil
	}

	return match, err
}

func GetNumberOfMatchByUserId(ctx context.Context, idUser int) (int, error) {
	var numberOfMatch = 0

	sql := `
		SELECT
			COUNT(*) AS number
		FROM ` + "`match`" + `
		WHERE id_user = ?
	`
	row := db.QueryRowContext(ctx, sql, idUser)
	if row.Err() != nil {
		return numberOfMatch, row.Err()
	}
	err := row.Scan(&numberOfMatch)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return numberOfMatch, nil
	}

	return numberOfMatch, err
}

func CreateMatch(ctx context.Context, m MatchCreate) (int, error) {
	sb := sqlbuilder.NewInsertBuilder()
	sb.InsertInto("`match`").Cols("id_user", "name").Values(m.IdUser, m.Name)
	sql, args := sb.Build()

	response, err := db.ExecContext(ctx, sql, args...)

	if err != nil {
		return 0, err
	}

	id, _ := response.LastInsertId()

	return int(id), err
}

func UpdateMatch(ctx context.Context, idMatch int, match MatchUpdate) error {
	canUpdate := false
	sb := sqlbuilder.NewUpdateBuilder()
	sb.Update("`match`").Where(sb.Equal("id", idMatch))

	if match.IdUser != nil {
		sb.SetMore(sb.Assign("id_user", *match.IdUser))
		canUpdate = true
	}

	if match.Name != nil {
		sb.SetMore(sb.Assign("name", *match.Name))
		canUpdate = true
	}

	if match.Round != nil {
		sb.SetMore(sb.Assign("round", *match.Round))
		canUpdate = true
	}

	if match.MultipleMasteryEnabled != nil {
		sb.SetMore(sb.Assign("multiple_mastery_enabled", *match.MultipleMasteryEnabled))
		canUpdate = true
	}
	if canUpdate {
		sql, args := sb.Build()
		_, err := db.ExecContext(ctx, sql, args...)
		return err
	}
	return errors.New("no parameters to update Match")
}

func IncreaseMatchRound(ctx context.Context, idMatch int) error {
	sql := "UPDATE `match` SET round = round + 1 WHERE id = ?"

	_, err := db.ExecContext(ctx, sql, idMatch)

	return err
}

func ResetMatch(ctx context.Context, idMatch int) error {
	var round = 0
	var multipleMasteryEnabled = 1
	err := UpdateMatch(ctx, idMatch, MatchUpdate{Round: &round, MultipleMasteryEnabled: &multipleMasteryEnabled})
	if err != nil {
		return err
	}

	err = ResetMatchPlayersSpells(ctx, idMatch)
	if err != nil {
		return err
	}

	return err
}

func DeleteMatch(ctx context.Context, idMatch int) error {
	sql := "DELETE FROM `match` WHERE id = ?"

	_, err := db.ExecContext(ctx, sql, idMatch)

	return err
}
