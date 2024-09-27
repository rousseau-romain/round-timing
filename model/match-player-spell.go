package model

import (
	"log"
	"strings"

	"github.com/rousseau-romain/round-timing/helper"
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

func GetSpellPlayerByIdSpellsPlayers(idSpellPlayer int) (MatchPlayerSpell, error) {
	matchSpell := MatchPlayerSpell{}

	sql := `
		SELECT
			mps.id,
			mps.match_id,
			mps.player_id,
			s.id,
			s.name,
			s.delay,
			` + helper.GetUrlImageSpellClause("s.id") + ` AS url_image,
			mps.round_before_recovery,
			mps.created_at,
			mps.updated_at
		FROM match_player_spell AS mps
		JOIN spell AS s ON s.id = mps.spell_id
		WHERE mps.id = ?
	`

	rows := db.QueryRow(sql, idSpellPlayer)

	if rows.Err() != nil {
		return matchSpell, rows.Err()
	}

	err := rows.Scan(&matchSpell.Id,
		&matchSpell.MatchId,
		&matchSpell.PlayerId,
		&matchSpell.Spell.Id,
		&matchSpell.Spell.Name,
		&matchSpell.Spell.Delay,
		&matchSpell.Spell.UrlImage,
		&matchSpell.RoundBeforeRecovery,
		&matchSpell.CreatedAt,
		&matchSpell.UpdatedAt,
	)

	return matchSpell, err
}

func GetSpellsPlayersByIdMatch(idMatch int) ([]MatchPlayerSpell, error) {
	// maitrise marteau id 138
	masteryIdSpells := []string{"134", "135", "136", "137", "139", "140", "141", "142"}
	matchSpells := []MatchPlayerSpell{}

	sql := `
		SELECT
			mps.id,
			mps.match_id,
			mps.player_id,
			s.id,
			s.name,
			s.delay,
			` + helper.GetUrlImageSpellClause("s.id") + ` AS url_image,
			mps.round_before_recovery,
			mps.created_at,
			mps.updated_at
		FROM match_player_spell AS mps
		JOIN spell AS s ON s.id = mps.spell_id
		JOIN ` + "`match`" + ` AS m ON m.id = mps.match_id 
		WHERE match_id = ?
		AND (m.multiple_mastery_enabled = 0 AND (s.id NOT IN (` + strings.Join(masteryIdSpells, ",") + `)))
		OR (m.multiple_mastery_enabled = 1)
	`
	// why can't use ?
	rows, err := db.Query(sql, idMatch)

	if err != nil {
		log.Println(err)
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
			&matchSpell.Spell.Delay,
			&matchSpell.Spell.UrlImage,
			&matchSpell.RoundBeforeRecovery,
			&matchSpell.CreatedAt,
			&matchSpell.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
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

	if err != nil {
		log.Println(err)
	}

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
