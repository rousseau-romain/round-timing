package model

type FeatureFlag struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func GetFeatureFlagIsEnabled(name string) bool {
	sql := `
		SELECT
			enabled
		FROM feature_flag
		WHERE name = "?"
	`

	row := db.QueryRow(sql, name)

	var isEnabled bool

	if row.Err() != nil {
		return false
	}
	err := row.Scan(&isEnabled)
	if err != nil {
		return false
	}

	return isEnabled
}
