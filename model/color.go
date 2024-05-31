package model

import (
	"log"
)

type Color struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetColorsTeam() ([]Color, error) {
	sql := `
		SELECT
			id,
			name
		FROM color_team
	`

	rows, err := db.Query(sql)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var colors []Color

	for rows.Next() {
		var color Color
		err := rows.Scan(
			&color.Id,
			&color.Name,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		colors = append(colors, color)
	}

	return colors, err
}
