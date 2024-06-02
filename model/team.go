package model

import (
	"log"
)

type Team struct {
	Id          int    `json:"id"`
	IdMatch     int    `json:"id_match"`
	IdColorTeam int    `json:"id_color_team"`
	Name        string `json:"name"`
	Color       string `json:"color"`
}

type TeamCreate struct {
	IdMatch     int
	IdColorTeam int
	Name        string
}

type TeamUpdate struct {
	IdColorTeam *int
	Name        *string
}

func GetTeamsByIdMatch(idMatch int) ([]Team, error) {
	sql := `
		SELECT
			t.id,
			t.id_match,
			t.id_color_team,
			t.name,
			ct.name AS color
		FROM team AS t
		JOIN color_team AS ct ON t.id_color_team = ct.id
		JOIN ` + "`match`" + ` AS m ON t.id_match = m.id
		JOIN user AS u ON m.id_user = u.id
		WHERE m.id = ?
	`

	rows, err := db.Query(sql, idMatch)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var teams []Team

	for rows.Next() {
		var team Team
		err := rows.Scan(
			&team.Id,
			&team.IdMatch,
			&team.IdColorTeam,
			&team.Name,
			&team.Color,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, err
}

func CreateTeam(m TeamCreate) (int, error) {

	sql := `
		INSERT INTO team (id_match, id_color_team, name)
		VALUES (?, ?, ?)
	`

	response, err := db.Exec(sql, m.IdMatch, m.IdColorTeam, m.Name)

	if err != nil {
		return 0, err
	}

	id, _ := response.LastInsertId()

	return int(id), err
}

func DeleteTeamsByMatchId(idMatch int) error {
	sql := `
		DELETE t FROM team AS t
		JOIN ` + "`match`" + ` AS m ON t.id_match = m.id
		WHERE m.id = ?
	`

	_, err := db.Exec(sql, idMatch)

	return err
}
