package user

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
