package model

import (
	"log"
	"round-timing/helper"
)

type Class struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UrlImage string `json:"url_image"`
}

func GetClasses() ([]Class, error) {
	sql := `
		SELECT
			id,
			name,
			` + helper.GetUrlImageClassClause("id") + ` AS url_image
		FROM class
		WHERE id != 13
	`

	rows, err := db.Query(sql)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var classes []Class

	for rows.Next() {
		var class Class
		err := rows.Scan(
			&class.Id,
			&class.Name,
			&class.UrlImage,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, err
}
