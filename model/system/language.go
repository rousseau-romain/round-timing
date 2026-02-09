package system

import "context"

type Language struct {
	Id     int    `json:"id"`
	Locale string `json:"locale"`
}

func GetLanguages(ctx context.Context) ([]Language, error) {
	sql := `
		SELECT
			id,
			locale
		FROM language
	`

	var languages []Language

	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var language Language
		err := rows.Scan(
			&language.Id,
			&language.Locale,
		)
		if err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}

func GetLanguagesIdByCode(ctx context.Context, code string) (int, error) {
	sql := `
		SELECT
			id
		FROM language
		WHERE locale = ?
	`

	row := db.QueryRowContext(ctx, sql, code)

	id := 0

	if row.Err() != nil {
		return id, row.Err()
	}
	err := row.Scan(&id)

	return id, err
}

func GetLanguageLocaleById(ctx context.Context, id int) (string, error) {
	sql := `
		SELECT
			locale
		FROM language
		WHERE id = ?
	`

	row := db.QueryRowContext(ctx, sql, id)

	var locale string

	if row.Err() != nil {
		return locale, row.Err()
	}
	err := row.Scan(&locale)

	return locale, err

}
