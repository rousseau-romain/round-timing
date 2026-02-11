package game

import (
	"context"

	"github.com/rousseau-romain/round-timing/pkg/sqlhelper"
)

type Spell struct {
	Id             int    `json:"id"`
	UrlImage       string `json:"url_image"`
	IdClass        int    `json:"id_class"`
	Name           string `json:"name"`
	ShortName      string `json:"short_name"`
	Color          string `json:"color"`
	Delay          int    `json:"delay"`
	IsGlobal       bool   `json:"is_global"`
	IsTeam         bool   `json:"is_team"`
	IsSelf         bool   `json:"is_self"`
	IsEndingCaster bool   `json:"is_ending_caster"`
}

func GetSpellsByIdClass(ctx context.Context, idLanguage int, idClass []int, idsToExclude []int) ([]Spell, error) {
	classPlaceholder, classArgs := sqlhelper.InClause(idClass)

	sql := `
		SELECT
			s.id,
			` + sqlhelper.URLImageSpellClause("s.id") + ` AS url_image,
			s.id_class,
			st.name,
			st.short_name,
			s.color,
			s.delay,
			s.is_global,
			s.is_team,
			s.is_self,
			s.is_ending_caster
		FROM spell s
		JOIN spell_translation st ON st.id_spell = s.id AND st.id_language = ?
		WHERE id_class IN ` + classPlaceholder + `
	`

	args := []interface{}{idLanguage}
	args = append(args, classArgs...)

	if len(idsToExclude) > 0 {
		excludePlaceholder, excludeArgs := sqlhelper.InClause(idsToExclude)
		sql = sql + "AND s.id NOT IN " + excludePlaceholder
		args = append(args, excludeArgs...)
	}
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []Spell

	for rows.Next() {
		var spell Spell
		err := rows.Scan(
			&spell.Id,
			&spell.UrlImage,
			&spell.IdClass,
			&spell.Name,
			&spell.ShortName,
			&spell.Color,
			&spell.Delay,
			&spell.IsGlobal,
			&spell.IsTeam,
			&spell.IsSelf,
			&spell.IsEndingCaster,
		)
		if err != nil {
			return nil, err
		}
		spells = append(spells, spell)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spells, nil
}
