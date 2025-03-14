package model

type UserConfiguration struct {
	Id              int    `json:"id"`
	IdUser          int    `json:"id_user"`
	IdConfiguration int    `json:"id_configuration"`
	Name            string `json:"name"`
	IsEnabled       bool   `json:"is_enabled"`
}

func GetConfigurationByIdConfigurationIdUser(idLanguage, idUser, idConfiguration int) (UserConfiguration, error) {
	sql := `
		SELECT
			IFNULL(uc.id, 0) AS id,
			IFNULL(uc.id_user, 0) AS id_user,
			IFNULL(c.id, 0) AS id_configuration,
			ct.name,
			IF(uc.id_user IS NULL, 0, 1) AS is_enabled
		FROM configuration AS c
		LEFT JOIN user_configuration AS uc ON uc.id_configuration = c.id
		JOIN configuration_translation AS ct ON ct.id_configuration = c.id AND ct.id_language = ?
		WHERE (uc.id_user = ? OR uc.id_user IS NULL) AND c.id = ?
	`

	rows := db.QueryRow(sql, idLanguage, idUser, idConfiguration)

	var userConfiguration UserConfiguration

	if rows.Err() != nil {
		return userConfiguration, rows.Err()
	}

	err := rows.Scan(
		&userConfiguration.Id,
		&userConfiguration.IdUser,
		&userConfiguration.IdConfiguration,
		&userConfiguration.Name,
		&userConfiguration.IsEnabled,
	)

	return userConfiguration, err
}

func GetAllConfigurationByIdUser(idLanguage, idUser int) ([]UserConfiguration, error) {
	sql := `
		SELECT
			IFNULL(uc.id, 0) AS id,
			IFNULL(uc.id_user, 0) AS id_user,
			IFNULL(c.id, 0) AS id_configuration,
			ct.name,
			IF(uc.id_user IS NULL, 0, 1) AS is_enabled
		FROM user_configuration AS uc
		RIGHT JOIN configuration AS c ON c.id = uc.id_configuration
		JOIN configuration_translation AS ct ON ct.id_configuration = c.id AND ct.id_language = ? 
		WHERE uc.id_user = ? OR uc.id_user IS NULL
	`

	rows, err := db.Query(sql, idLanguage, idUser)

	if err != nil {
		return nil, err
	}

	var configurationByIdUser []UserConfiguration

	for rows.Next() {
		var configuration UserConfiguration
		err := rows.Scan(
			&configuration.Id,
			&configuration.IdUser,
			&configuration.IdConfiguration,
			&configuration.Name,
			&configuration.IsEnabled,
		)
		if err != nil {
			return nil, err
		}
		configurationByIdUser = append(configurationByIdUser, configuration)
	}

	return configurationByIdUser, err
}

func ToggleUserConfiguration(idUser, idConfiguration int) error {
	row := db.QueryRow("SELECT EXISTS (SELECT id FROM user_configuration WHERE id_user = ? AND id_configuration = ?)", idUser, idConfiguration)

	var isEnable bool

	if row.Err() != nil {
		return row.Err()
	}

	err := row.Scan(&isEnable)

	if err != nil {
		return err
	}

	if isEnable {
		_, err := db.Exec("DELETE FROM user_configuration WHERE id_user = ? AND id_configuration = ?", idUser, idConfiguration)
		return err
	} else {
		_, err := db.Exec("INSERT INTO user_configuration (id_user, id_configuration) VALUES (?, ?)", idUser, idConfiguration)
		return err
	}
}
