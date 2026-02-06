package user

import "context"

func IsEmailWhiteListed(ctx context.Context, email string) (bool, error) {
	var count int
	sql := `
		SELECT
			COUNT(*)
		FROM email_white_listed
		WHERE email = ?;
	`
	err := db.QueryRowContext(ctx, sql, email).Scan(&count)

	if count > 0 {
		return true, err
	}

	return false, err
}
