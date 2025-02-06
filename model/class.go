package model

import (
	"log"

	"github.com/rousseau-romain/round-timing/helper"
)

type Class struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UrlImage string `json:"url_image"`
}

func GetClasses(idLanguage int) ([]Class, error) {
	sql := `
		SELECT
			c.id,
			cn.name,
			? AS url_image
		FROM class c
		JOIN class_translation cn ON cn.id_class = c.id AND cn.id_language = ?
		WHERE c.id != 13
	`

	rows, err := db.Query(sql, helper.GetUrlImageClassClause("c.id"), idLanguage)

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
