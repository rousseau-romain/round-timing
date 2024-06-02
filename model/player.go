package model

import (
	"errors"
	"log"
	"round-timing/helper"

	"github.com/huandu/go-sqlbuilder"
)

type Player struct {
	Id     int `json:"id"`
	Idteam int `json:"id_team"`
	// IdClass int    `json:"id_class"`
	Name  string `json:"name"`
	Class Class  `json:"class"`
	Team  Team   `json:"team"`
}

type PlayerCreate struct {
	IdTeam  int
	IdClass int
	Name    string
}

type PlayerUpdate struct {
	IdClass *int
	Name    *string
}

func GetPlayersByIdMatch(idMatch int) ([]Player, error) {
	sql := `
		SELECT
			p.id,
			p.id_team,
			p.name,
			p.id_class,
			c.name AS class_name,
			` + helper.GetUrlImageClassClause("c.id") + ` AS url_image
			t.id AS id_team,
			t.name AS team_name,
			ct.name AS color_team
		FROM player AS p
		JOIN class AS c ON c.id = p.id_class
		JOIN team AS t ON t.id = p.id_team
		JOIN color_team AS ct ON ct.id = t.id_color_team
		JOIN ` + "`match`" + ` AS m ON m.id = t.id_match
		WHERE m.id = ? 
	`

	rows, err := db.Query(sql, idMatch)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var players []Player

	for rows.Next() {
		var player Player
		err := rows.Scan(
			&player.Id,
			&player.Idteam,
			&player.Name,
			&player.Class.Id,
			&player.Class.Name,
			&player.Class.UrlImage,
			&player.Team.Id,
			&player.Team.Name,
			&player.Team.Color,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		players = append(players, player)
	}

	return players, err
}

func GetPlayer(idPlayer int) (Player, error) {
	sql := `
		SELECT 
			p.id,
			p.id_team,
			p.name,
			p.id_class,
			c.name AS class_name,
			` + helper.GetUrlImageClassClause("c.id") + ` AS url_image
			t.id AS id_team,
			t.name AS team_name,
			ct.name AS color_team
		FROM player AS p
		JOIN team AS t ON t.id = p.id_team
		JOIN color_team AS ct ON ct.id = t.id_color_team
		JOIN class AS c ON c.id = p.id_class
		WHERE p.id = ?
	`

	var player Player
	err := db.QueryRow(sql, idPlayer).Scan(
		&player.Id,
		&player.Idteam,
		&player.Name,
		&player.Class.Id,
		&player.Class.Name,
		&player.Class.UrlImage,
		&player.Team.Id,
		&player.Team.Name,
		&player.Team.Color,
	)

	if err != nil {
		log.Println(err)
		return Player{}, err
	}

	return player, err
}

func CreatePlayer(p PlayerCreate) (int, error) {

	sql := `
		INSERT INTO player (id_class, id_team, name)
		VALUES (?, ?, ?)
	`

	response, err := db.Exec(sql, p.IdClass, p.IdTeam, p.Name)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	id, _ := response.LastInsertId()

	log.Println("zefsdfsdfsd", id)

	return int(id), err
}

func UpdatePlayer(idPlayer int, player PlayerUpdate) error {
	canUpdate := false
	sb := sqlbuilder.NewUpdateBuilder()
	sb.Update("player").Where(sb.Equal("id", idPlayer))

	if player.IdClass != nil {
		sb.Set(sb.Assign("id_class", *player.IdClass))
		canUpdate = true
	}

	if player.Name != nil {
		sb.Set(sb.Assign("name", *player.Name))
		canUpdate = true
	}

	if canUpdate {
		sql, args := sb.Build()
		_, err := db.Exec(sql, args...)
		return err
	}

	return errors.New("no parameters to update Player")
}

func DeletePlayersByMatchId(idMatch int) error {
	sql := `
		DELETE p FROM player AS p
		JOIN team AS t ON p.id_team = t.id
		JOIN ` + "`match`" + ` AS m ON t.id_match = m.id
		WHERE m.id = ?
	`

	_, err := db.Exec(sql, idMatch)

	return err
}

func DeletePlayer(idPlayer int) error {
	sql := "DELETE FROM player WHERE id = ?"
	_, err := db.Exec(sql, idPlayer)

	return err
}
