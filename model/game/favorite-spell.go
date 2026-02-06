package game

import (
	"context"

	"github.com/rousseau-romain/round-timing/pkg/sqlhelper"
)

type SpellByClass struct {
	IdSpell    int    `json:"id_spell"`
	IdClass    int    `json:"id_class"`
	UrlImage   string `json:"url_image"`
	Name       string `json:"name"`
	IsFavorite bool   `json:"is_favorite"`
}

func GetFavoriteSpellByIdUserAndIdSpell(ctx context.Context, idLanguage, idUser, idSpell int) (SpellByClass, error) {
	sql := `
		SELECT
			s.id AS id_spell,
			s.id_class AS id_class,
			` + sqlhelper.URLImageSpellClause("s.id") + ` AS spell_url_image,
			st.name AS spell_name,
			IF(fs.id_spell IS NULL , 0, 1) AS is_favorite
		FROM spell AS s
		JOIN spell_translation AS st ON st.id_spell = s.id AND st.id_language = ?
		JOIN class AS c ON s.id_class = c.id
		LEFT JOIN favorite_spell fs ON fs.id_spell = s.id AND fs.id_user = ?
		WHERE s.id = ?
	`

	row := db.QueryRowContext(ctx, sql, idLanguage, idUser, idSpell)

	var spellByClass SpellByClass

	if row.Err() != nil {
		return spellByClass, row.Err()
	}

	err := row.Scan(
		&spellByClass.IdSpell,
		&spellByClass.IdClass,
		&spellByClass.UrlImage,
		&spellByClass.Name,
		&spellByClass.IsFavorite,
	)

	return spellByClass, err
}

func GetFavoriteSpellsByIdUser(ctx context.Context, idLanguage, idUser int) ([]SpellByClass, error) {
	sql := `
		SELECT
			s.id AS id_spell,
			s.id_class AS id_class,
			` + sqlhelper.URLImageSpellClause("s.id") + ` AS spell_url_image,
			st.name AS spell_name,
			IF(fs.id_spell IS NULL , 0, 1) AS is_favorite
		FROM spell AS s
		JOIN spell_translation AS st ON st.id_spell = s.id AND st.id_language = ?
		JOIN class AS c ON s.id_class = c.id
		LEFT JOIN favorite_spell fs ON fs.id_spell = s.id AND fs.id_user = ?
		ORDER BY is_favorite DESC, fs.id ASC
	`

	rows, err := db.QueryContext(ctx, sql, idLanguage, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spellByClasses []SpellByClass

	for rows.Next() {
		var spellByClass SpellByClass
		err := rows.Scan(
			&spellByClass.IdSpell,
			&spellByClass.IdClass,
			&spellByClass.UrlImage,
			&spellByClass.Name,
			&spellByClass.IsFavorite,
		)
		if err != nil {
			return nil, err
		}
		spellByClasses = append(spellByClasses, spellByClass)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spellByClasses, nil
}

func ToggleIsFavoriteSpell(ctx context.Context, idUser, idSpell int) error {
	isFavorite := false

	rows, err := db.QueryContext(ctx, "SELECT id FROM favorite_spell WHERE id_user = ? AND id_spell = ?", idUser, idSpell)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		isFavorite = true
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if isFavorite {
		_, err := db.ExecContext(ctx, "DELETE FROM favorite_spell WHERE id_user = ? AND id_spell = ?", idUser, idSpell)
		return err
	} else {
		_, err := db.ExecContext(ctx, "INSERT INTO favorite_spell (id_user, id_spell) VALUES (?, ?)", idUser, idSpell)
		return err
	}
}
