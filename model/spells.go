package model

import (
	"log"
	"strconv"
	"strings"

	"github.com/rousseau-romain/round-timing/helper"
)

type Spell struct {
	Id             int    `json:"id"`
	UrlImage       string `json:"url_image"`
	IdClass        int    `json:"id_class"`
	Name           string `json:"name"`
	Delay          int    `json:"delay"`
	IsGlobal       bool   `json:"is_global"`
	IsTeam         bool   `json:"is_team"`
	IsSelf         bool   `json:"is_self"`
	IsEndingCaster bool   `json:"is_ending_caster"`
}

func GetSpellsByIdCLass(idClass []int) ([]Spell, error) {
	var strIdClass []string
	for _, id := range idClass {
		strIdClass = append(strIdClass, strconv.Itoa(id))
	}
	sql := `
		SELECT
			id,
			` + helper.GetUrlImageSpellClause("id") + ` AS url_image,
			id_class,
			name,
			delay,
			is_global,
			is_team,
			is_self,
			is_ending_caster
		FROM spell
		WHERE id_class IN (` + strings.Join(strIdClass, ",") + `)
	`

	rows, err := db.Query(sql)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var spells []Spell

	for rows.Next() {
		var spell Spell
		err := rows.Scan(
			&spell.Id,
			&spell.UrlImage,
			&spell.IdClass,
			&spell.Name,
			&spell.Delay,
			&spell.IsGlobal,
			&spell.IsTeam,
			&spell.IsSelf,
			&spell.IsEndingCaster,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		spells = append(spells, spell)
	}

	return spells, err
}
