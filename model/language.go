package model

import (
	"log"
)

type Language struct {
	Id     int    `json:"id"`
	Locale string `json:"locale"`
}

func GetLanguages() ([]Language, error) {
	sql := `
		SELECT
			id,
			locale
		FROM language
	`

	var languages []Language

	rows, err := db.Query(sql)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var language Language
		err := rows.Scan(
			&language.Id,
			&language.Locale,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		languages = append(languages, language)
	}

	return languages, err
}

func GetLanguagesIdByCode(code string) (int, error) {
	sql := `
		SELECT
			id
		FROM language
		WHERE locale = ?
	`

	row := db.QueryRow(sql, code)

	id := 0

	if row.Err() != nil {
		return id, row.Err()
	}
	err := row.Scan(&id)

	return id, err
}

func GetLanguageLocaleById(id int) (string, error) {
	sql := `
		SELECT
			locale
		FROM language
		WHERE id = ?
	`

	row := db.QueryRow(sql, id)

	var locale string

	if row.Err() != nil {
		return locale, row.Err()
	}
	err := row.Scan(&locale)

	return locale, err

}
