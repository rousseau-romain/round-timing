package model

import (
	"fmt"

	"github.com/rousseau-romain/round-timing/pkg/constants"
	"github.com/rousseau-romain/round-timing/pkg/sqlhelper"
)

type MatchPlayerSpell struct {
	Id                  int    `json:"id"`
	MatchId             int    `json:"match_id"`
	PlayerId            int    `json:"player_id"`
	Spell               Spell  `json:"spell"`
	RoundBeforeRecovery int    `json:"round_before_recovery"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type MatchPlayerSpellCreate struct {
	MatchId             int
	PlayerId            int
	SpellId             int
	RoundBeforeRecovery int
}

type MatchPlayerSpellUpdate struct {
	RoundBeforeRecovery *string
}

func GetSpellPlayerByIdSpellsPlayers(idLanguage int, idSpellPlayer int) (MatchPlayerSpell, error) {
	matchSpell := MatchPlayerSpell{}

	sql := `
		SELECT
			mps.id,
			mps.match_id,
			mps.player_id,
			s.id,
			st.name,
			st.short_name,
			s.color,
			s.delay,
			` + sqlhelper.URLImageSpellClause("s.id") + ` AS url_image,
			mps.round_before_recovery,
			mps.created_at,
			mps.updated_at
		FROM match_player_spell AS mps
		JOIN spell AS s ON s.id = mps.spell_id
		JOIN spell_translation AS st ON st.id_spell = s.id	 AND st.id_language = ?
		WHERE mps.id = ?
	`

	rows := db.QueryRow(sql, idLanguage, idSpellPlayer)

	if rows.Err() != nil {
		return matchSpell, rows.Err()
	}

	err := rows.Scan(&matchSpell.Id,
		&matchSpell.MatchId,
		&matchSpell.PlayerId,
		&matchSpell.Spell.Id,
		&matchSpell.Spell.Name,
		&matchSpell.Spell.ShortName,
		&matchSpell.Spell.Color,
		&matchSpell.Spell.Delay,
		&matchSpell.Spell.UrlImage,
		&matchSpell.RoundBeforeRecovery,
		&matchSpell.CreatedAt,
		&matchSpell.UpdatedAt,
	)

	return matchSpell, err
}

func GetSpellsPlayersByIdMatch(idLanguage, idMatch, idUser int, getOnlyFavorite bool) ([]MatchPlayerSpell, error) {
	matchSpells := []MatchPlayerSpell{}

	joinFavoriteClause := "LEFT"
	if getOnlyFavorite {
		joinFavoriteClause = "INNER"

	}

	masteryClause := "m.multiple_mastery_enabled = 0 AND (s.id NOT IN ("
	for _, id := range constants.MasteryIdSpells {
		masteryClause += fmt.Sprintf("%d, ", id)
	}
	masteryClause = masteryClause[:len(masteryClause)-2]
	masteryClause += "))"

	sql := `
		SELECT
			mps.id,
			mps.match_id,
			mps.player_id,
			s.id,
			st.name,
			st.short_name,
			s.color,
			s.delay,
			` + sqlhelper.URLImageSpellClause("s.id") + ` AS url_image,
			mps.round_before_recovery,
			mps.created_at,
			mps.updated_at
		FROM match_player_spell AS mps
		JOIN spell AS s ON s.id = mps.spell_id
		JOIN spell_translation AS st ON st.id_spell = s.id AND st.id_language = ?
		JOIN ` + "`match`" + ` AS m ON m.id = mps.match_id
		` + joinFavoriteClause + ` JOIN favorite_spell fs ON fs.id_spell = s.id AND fs.id_user = ?
		WHERE (
			mps.match_id = ?
			AND (
				(
					` + masteryClause + `
				)
				OR (m.multiple_mastery_enabled = 1)
			)
		)
		ORDER BY mps.player_id DESC, IF(fs.id_spell IS NULL , 1, 2) DESC, fs.id ASC ,s.id ASC
	`
	// why can't use ?
	rows, err := db.Query(sql, idLanguage, idUser, idMatch)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var matchSpell MatchPlayerSpell
		err := rows.Scan(
			&matchSpell.Id,
			&matchSpell.MatchId,
			&matchSpell.PlayerId,
			&matchSpell.Spell.Id,
			&matchSpell.Spell.Name,
			&matchSpell.Spell.ShortName,
			&matchSpell.Spell.Color,
			&matchSpell.Spell.Delay,
			&matchSpell.Spell.UrlImage,
			&matchSpell.RoundBeforeRecovery,
			&matchSpell.CreatedAt,
			&matchSpell.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		matchSpells = append(matchSpells, matchSpell)
	}

	return matchSpells, err
}

func CreateMatchPlayersSpells(matchPlayersSpells []MatchPlayerSpellCreate) error {
	sql := `
		INSERT INTO match_player_spell
			(match_id, player_id, spell_id, round_before_recovery)
		VALUES
			
	`
	args := []interface{}{}
	for _, mps := range matchPlayersSpells {
		sql += "(?, ?, ?, ?),"
		args = append(args, mps.MatchId, mps.PlayerId, mps.SpellId, mps.RoundBeforeRecovery)
	}
	sql = sql[0 : len(sql)-1]

	_, err := db.Exec(sql, args...)

	return err
}

func DecreasePlayersSpellsRoundBeforeRecoveryByIdMatch(idMatch int) error {
	sql := `
		UPDATE match_player_spell
		SET round_before_recovery = CASE
			WHEN round_before_recovery > 0 THEN round_before_recovery - 1
			ELSE 0
		END
		WHERE match_id = ?
	`

	_, err := db.Exec(sql, idMatch)

	return err
}

func UsePlayerSpellByIdPlayerSpell(idPlayerSpell int) error {
	sql := `
		UPDATE match_player_spell 
		SET round_before_recovery = (
			SELECT s.delay  
			FROM match_player_spell AS mps 
			JOIN spell AS s ON s.id = mps.spell_id 
			WHERE mps.id = ?
		)
		WHERE id = ?
	`
	_, err := db.Exec(sql, idPlayerSpell, idPlayerSpell)

	return err
}

func RemoveRoundRecoverySpellByIdPlayerSpell(idPlayerSpell int) error {
	sql := `
		UPDATE match_player_spell 
		SET round_before_recovery = CASE
			WHEN round_before_recovery > 0 THEN round_before_recovery - 1
			ELSE 0
		END
		WHERE id = ?
	`
	_, err := db.Exec(sql, idPlayerSpell)

	return err
}

func ResetMatchPlayersSpells(idMatch int) error {
	sql := `
		DELETE FROM match_player_spell
		WHERE match_id = ?
	`

	_, err := db.Exec(sql, idMatch)

	return err

}

func DeleteMatchPlayersSpellsByMatchId(idMatch int) error {
	sql := "DELETE FROM match_player_spell WHERE match_id = ?"

	_, err := db.Exec(sql, idMatch)

	return err
}

func DeleteMatchPlayersSpellsByPlayer(idPlayer int) error {
	sql := "DELETE FROM match_player_spell WHERE player_id = ?"

	_, err := db.Exec(sql, idPlayer)

	return err
}
