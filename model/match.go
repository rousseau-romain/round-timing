package model

import (
	"errors"
	"log"

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

func GetMatchsByIdUser(idUser int) ([]Match, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "id_user", "name", "round", "multiple_mastery_enabled", "created_at", "updated_at").From("`match`").Where(sb.Equal("id_user", idUser))
	sql, args := sb.Build()

	rows, err := db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	var matchs []Match

	for rows.Next() {
		var match Match
		err := rows.Scan(&match.Id, &match.IdUser, &match.Name, &match.Round, &match.MultipleMasteryEnabled, &match.CreatedAt, &match.UpdatedAt)
		if err != nil {
			return matchs, err
		}
		matchs = append(matchs, match)
	}

	return matchs, err
}

func GetMatch(idMatch int) (Match, error) {
	match := Match{}

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "id_user", "name", "round", "multiple_mastery_enabled", "created_at", "updated_at").From("`match`").Where(sb.Equal("id", idMatch))
	sql, args := sb.Build()

	rows := db.QueryRow(sql, args...)
	if rows.Err() != nil {
		return match, rows.Err()
	}
	err := rows.Scan(&match.Id, &match.IdUser, &match.Name, &match.Round, &match.MultipleMasteryEnabled, &match.CreatedAt, &match.UpdatedAt)

	return match, err
}

func GetLastMatchByUserId(idUser int) (Match, error) {
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
	rows := db.QueryRow(sql, idUser, idUser)
	if rows.Err() != nil {
		return match, rows.Err()
	}
	err := rows.Scan(&match.Id, &match.IdUser, &match.Name, &match.Round, &match.MultipleMasteryEnabled, &match.CreatedAt, &match.UpdatedAt)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return match, nil
	}

	return match, err
}

func GetNumberOfMatchByUserId(idUser int) (int, error) {
	var numberOfMatch = 0

	sql := `
		SELECT
			COUNT(*) AS number
		FROM ` + "`match`" + `
		WHERE id_user = ?
	`
	rows := db.QueryRow(sql, idUser)
	if rows.Err() != nil {
		return numberOfMatch, rows.Err()
	}
	err := rows.Scan(&numberOfMatch)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return numberOfMatch, nil
	}

	return numberOfMatch, err
}

func CreateMatch(m MatchCreate) (int, error) {
	sb := sqlbuilder.NewInsertBuilder()
	sb.InsertInto("`match`").Cols("id_user", "name").Values(m.IdUser, m.Name)
	sql, args := sb.Build()

	response, err := db.Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	id, _ := response.LastInsertId()

	return int(id), err
}

func UpdateMatch(idMatch int, match MatchUpdate) error {
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
		_, err := db.Exec(sql, args...)
		return err
	}
	return errors.New("no parameters to update Match")
}

func IncreaseMatchRound(idMatch int) error {
	sql := "UPDATE `match` SET round = round + 1 WHERE id = ?"

	_, err := db.Exec(sql, idMatch)

	return err
}

func ResetMatch(idMatch int) error {
	var round = 0
	var multipleMasteryEnabled = 1
	err := UpdateMatch(idMatch, MatchUpdate{Round: &round, MultipleMasteryEnabled: &multipleMasteryEnabled})
	if err != nil {
		log.Println(err)
		return err
	}

	err = ResetMatchPlayersSpells(idMatch)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func DeleteMatch(idMatch int) error {
	sql := "DELETE FROM `match` WHERE id = ?"

	_, err := db.Exec(sql, idMatch)

	return err
}
