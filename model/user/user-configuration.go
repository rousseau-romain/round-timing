package user

import (
	"context"
	"strings"
)

type ConfigurationValueLabel struct {
	Value string
	Label string
}

type UserConfiguration struct {
	Id              int    `json:"id"`
	IdUser          int    `json:"id_user"`
	IdConfiguration int    `json:"id_configuration"`
	Key             string `json:"key"`
	Name            string `json:"name"`
	Value           string `json:"value"`
	DefaultValue    string `json:"default_value"`
	PossibleValues  string `json:"possible_values"`
	ValueLabelsRaw  string `json:"-"`
}

func (uc UserConfiguration) PossibleValuesList() []string {
	return strings.Split(uc.PossibleValues, ",")
}

func (uc UserConfiguration) IsBooleanConfig() bool {
	values := uc.PossibleValuesList()
	if len(values) != 2 {
		return false
	}
	hasTrue := false
	hasFalse := false
	for _, v := range values {
		if v == "true" {
			hasTrue = true
		}
		if v == "false" {
			hasFalse = true
		}
	}
	return hasTrue && hasFalse
}

// PossibleValuesWithLabels parses ValueLabelsRaw ("dark:Sombre|light:Clair|auto:Auto")
// and returns them in the order defined by PossibleValues.
func (uc UserConfiguration) PossibleValuesWithLabels() []ConfigurationValueLabel {
	labelMap := make(map[string]string)
	if uc.ValueLabelsRaw != "" {
		for _, pair := range strings.Split(uc.ValueLabelsRaw, "|") {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				labelMap[parts[0]] = parts[1]
			}
		}
	}

	var result []ConfigurationValueLabel
	for _, v := range uc.PossibleValuesList() {
		label := v
		if l, ok := labelMap[v]; ok {
			label = l
		}
		result = append(result, ConfigurationValueLabel{Value: v, Label: label})
	}
	return result
}

func (uc UserConfiguration) GetValueLabel(value string) string {
	for _, vl := range uc.PossibleValuesWithLabels() {
		if vl.Value == value {
			return vl.Label
		}
	}
	return value
}

func GetConfigurationByKeyAndIdUser(ctx context.Context, idLanguage, idUser int, key string) (UserConfiguration, error) {
	sql := `
		SELECT
			IFNULL(uc.id, 0) AS id,
			IFNULL(uc.id_user, 0) AS id_user,
			IFNULL(c.id, 0) AS id_configuration,
			c.` + "`key`" + `,
			ct.name,
			IFNULL(uc.value, c.default_value) AS value,
			c.default_value,
			c.possible_values,
			IFNULL((SELECT GROUP_CONCAT(CONCAT(cvt.value, ':', cvt.label) ORDER BY cvt.id SEPARATOR '|')
				FROM configuration_value_translation cvt
				WHERE cvt.id_configuration = c.id AND cvt.id_language = ?), '') AS value_labels
		FROM configuration AS c
		LEFT JOIN user_configuration AS uc ON uc.id_configuration = c.id AND uc.id_user = ?
		JOIN configuration_translation AS ct ON ct.id_configuration = c.id AND ct.id_language = ?
		WHERE c.key = ?
	`

	row := db.QueryRowContext(ctx, sql, idLanguage, idUser, idLanguage, key)

	var userConfiguration UserConfiguration

	if row.Err() != nil {
		return userConfiguration, row.Err()
	}

	err := row.Scan(
		&userConfiguration.Id,
		&userConfiguration.IdUser,
		&userConfiguration.IdConfiguration,
		&userConfiguration.Key,
		&userConfiguration.Name,
		&userConfiguration.Value,
		&userConfiguration.DefaultValue,
		&userConfiguration.PossibleValues,
		&userConfiguration.ValueLabelsRaw,
	)

	if err != nil && err.Error() != "sql: no rows in result set" {
		return userConfiguration, err
	}

	return userConfiguration, nil
}

func GetConfigurationByIdConfigurationIdUser(ctx context.Context, idLanguage, idUser, idConfiguration int) (UserConfiguration, error) {
	sql := `
		SELECT
			IFNULL(uc.id, 0) AS id,
			IFNULL(uc.id_user, 0) AS id_user,
			IFNULL(c.id, 0) AS id_configuration,
			c.` + "`key`" + `,
			ct.name,
			IFNULL(uc.value, c.default_value) AS value,
			c.default_value,
			c.possible_values,
			IFNULL((SELECT GROUP_CONCAT(CONCAT(cvt.value, ':', cvt.label) ORDER BY cvt.id SEPARATOR '|')
				FROM configuration_value_translation cvt
				WHERE cvt.id_configuration = c.id AND cvt.id_language = ?), '') AS value_labels
		FROM configuration AS c
		LEFT JOIN user_configuration AS uc ON uc.id_configuration = c.id AND uc.id_user = ?
		JOIN configuration_translation AS ct ON ct.id_configuration = c.id AND ct.id_language = ?
		WHERE c.id = ?
	`

	row := db.QueryRowContext(ctx, sql, idLanguage, idUser, idLanguage, idConfiguration)

	var userConfiguration UserConfiguration

	if row.Err() != nil {
		return userConfiguration, row.Err()
	}

	err := row.Scan(
		&userConfiguration.Id,
		&userConfiguration.IdUser,
		&userConfiguration.IdConfiguration,
		&userConfiguration.Key,
		&userConfiguration.Name,
		&userConfiguration.Value,
		&userConfiguration.DefaultValue,
		&userConfiguration.PossibleValues,
		&userConfiguration.ValueLabelsRaw,
	)

	if err != nil && err.Error() != "sql: no rows in result set" {
		return userConfiguration, err
	}

	return userConfiguration, nil
}

func GetAllConfigurationByIdUser(ctx context.Context, idLanguage, idUser int) ([]UserConfiguration, error) {
	sql := `
		SELECT
			IFNULL(uc.id, 0) AS id,
			IFNULL(uc.id_user, 0) AS id_user,
			c.id AS id_configuration,
			c.` + "`key`" + `,
			ct.name,
			IFNULL(uc.value, c.default_value) AS value,
			c.default_value,
			c.possible_values,
			IFNULL((SELECT GROUP_CONCAT(CONCAT(cvt.value, ':', cvt.label) ORDER BY cvt.id SEPARATOR '|')
				FROM configuration_value_translation cvt
				WHERE cvt.id_configuration = c.id AND cvt.id_language = ?), '') AS value_labels
		FROM configuration AS c
		LEFT JOIN user_configuration AS uc ON uc.id_configuration = c.id AND uc.id_user = ?
		JOIN configuration_translation AS ct ON ct.id_configuration = c.id AND ct.id_language = ?
	`

	rows, err := db.QueryContext(ctx, sql, idLanguage, idUser, idLanguage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configurationByIdUser []UserConfiguration

	for rows.Next() {
		var configuration UserConfiguration
		err := rows.Scan(
			&configuration.Id,
			&configuration.IdUser,
			&configuration.IdConfiguration,
			&configuration.Key,
			&configuration.Name,
			&configuration.Value,
			&configuration.DefaultValue,
			&configuration.PossibleValues,
			&configuration.ValueLabelsRaw,
		)
		if err != nil {
			return nil, err
		}
		configurationByIdUser = append(configurationByIdUser, configuration)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configurationByIdUser, nil
}

func SetUserConfiguration(ctx context.Context, idUser, idConfiguration int, value string) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO user_configuration (id_user, id_configuration, value) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE value = ?",
		idUser, idConfiguration, value, value,
	)
	return err
}

func ToggleUserConfiguration(ctx context.Context, idUser, idConfiguration, idLanguage int) error {
	// Get current effective value
	uc, err := GetConfigurationByIdConfigurationIdUser(ctx, idLanguage, idUser, idConfiguration)
	if err != nil {
		return err
	}

	newValue := "true"
	if uc.Value == "true" {
		newValue = "false"
	}

	return SetUserConfiguration(ctx, idUser, idConfiguration, newValue)
}
