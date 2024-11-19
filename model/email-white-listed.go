package model

import "log"

func IsEmailWhiteListed(email string) (bool, error) {
	var count int
	sql := `
		SELECT
			COUNT(*)
		FROM email_white_listed
		WHERE email = ?;
	`
	err := db.QueryRow(sql, email).Scan(&count)

	if count > 0 {
		return true, err
	}

	return false, err
}

func CreateEmailWhiteListed(email string) (int64, error) {
	sql := `INSERT INTO email_white_listed (email) VALUES (?)`

	response, err := db.Exec(sql, email)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	id, _ := response.LastInsertId()

	return id, err
}
