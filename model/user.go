package model

import (
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
)

type User struct {
	Id            int    `json:"id"`
	ProviderLogin string `json:"provider_login"`
	Oauth2Id      string `json:"oauth2_id"`
	Enabled       bool   `json:"enabled"`
	Email         string `json:"email"`
	Hash          string `json:"hash"`
	IdLanguage    int    `json:"id_language"`
	IsAdmin       bool   `json:"is_admin"`
	IdShare       string `json:"id_share"`
}
type UserUpdate struct {
	Enabled    *bool
	IdLanguage *int
}

type UserCreate struct {
	ProviderLogin string
	Oauth2Id      string
	Email         string
	Hash          string
	IdLanguage    int
}

func GetUsers() ([]User, error) {
	sql := `
		SELECT 
			id,
			email,
			id_language,
			enabled,
			is_admin,
			id_share
		FROM user
	`

	rows, err := db.Query(sql)

	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.IdLanguage,
			&user.Enabled,
			&user.IsAdmin,
			&user.IdShare,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}

func GetUserByEmail(email string) (User, error) {
	user := User{}

	sql := `
			SELECT
				id,
				oauth2_id,
				enabled,
				email,
				provider_login,
				hash,
				id_language,
				is_admin,
				id_share
			FROM user
			WHERE email = ?
		`

	err := db.QueryRow(sql, email).Scan(
		&user.Id,
		&user.Oauth2Id,
		&user.Enabled,
		&user.Email,
		&user.ProviderLogin,
		&user.Hash,
		&user.IdLanguage,
		&user.IsAdmin,
		&user.IdShare,
	)

	if err != nil {
		return user, err
	}

	return user, err
}

func GetUserById(idUser int) (User, error) {
	user := User{}

	sql := `
			SELECT
				id,
				oauth2_id,
				enabled,
				email,
				id_language,
				is_admin,
				id_share
			FROM user
			WHERE id = ?
		`

	err := db.QueryRow(sql, idUser).Scan(
		&user.Id,
		&user.Oauth2Id,
		&user.Enabled,
		&user.Email,
		&user.IdLanguage,
		&user.IsAdmin,
		&user.IdShare,
	)

	if err != nil {
		return user, err
	}

	return user, err
}

func GetUserIdByMatch(idmatch int) (User, error) {
	user := User{}

	sql := fmt.Sprintf(`
		SELECT
			u.id,
			u.oauth2_id,
			u.enabled,
			u.email,
			u.id_language,
			u.is_admin,
			u.id_share
		FROM %s AS m
		JOIN user AS u ON m.id_user = u.id
		WHERE m.id = ?
	`, "`match`")

	err := db.QueryRow(sql, idmatch).Scan(&user.Id, &user.Oauth2Id, &user.Enabled, &user.Email, &user.IdLanguage, &user.IsAdmin, &user.IdShare)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user, nil
	}

	return user, err
}

func GetUserByOauth2Id(oauth2Id string) (User, error) {
	user := User{}

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "oauth2_id", "enabled", "email", "provider_login", "id_language", "is_admin", "id_share").From("user").Where(sb.Equal("oauth2_id", oauth2Id))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&user.Id, &user.Oauth2Id, &user.Enabled, &user.Email, &user.ProviderLogin, &user.IdLanguage, &user.IsAdmin, &user.IdShare)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user, nil
	}
	return user, err
}

func UserExistsByOauth2Id(oauth2Id string) (bool, error) {
	userId := 0

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id").From("user").Where(sb.Equal("oauth2_id", oauth2Id))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&userId)

	if err != nil {
		return false, nil
	}

	return userId != 0, err
}

func UserExistsByIdShare(idShare string) (bool, error) {
	userId := 0

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id").From("user").Where(sb.Equal("id_share", idShare))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&userId)

	if err != nil {
		return false, nil
	}

	return userId != 0, err
}

func UserExistsByEmail(email string) (string, error) {
	var providerLoginName string

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("provider_login").From("user").Where(sb.Equal("email", email))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&providerLoginName)

	if err != nil {
		return "", nil
	}

	return providerLoginName, err
}

func IsAdminUser(userId int) (bool, error) {
	var isAdmin int
	sql := `
		SELECT
			is_admin
		FROM user
		WHERE user = ?;
	`
	err := db.QueryRow(sql, userId).Scan(&isAdmin)

	if isAdmin == 1 {
		return true, err
	}

	return false, err
}

func CreateUser(user UserCreate) (int64, error) {

	sql := `
		INSERT INTO user (provider_login, oauth2_id, email, id_language, hash)
		VALUES (?, ?, ?, ?, ?)
	`

	response, err := db.Exec(sql, user.ProviderLogin, user.Oauth2Id, user.Email, user.IdLanguage, user.Hash)

	if err != nil {
		return 0, err
	}

	id, _ := response.LastInsertId()

	return id, err
}

func UpdateUser(idUser int, user UserUpdate) error {
	canUpdate := false
	sb := sqlbuilder.NewUpdateBuilder()
	sb.Update("user").Where(sb.Equal("id", idUser))

	if user.Enabled != nil {
		sb.Set(sb.Assign("enabled", *user.Enabled))
		canUpdate = true
	}

	if user.IdLanguage != nil {
		sb.Set(sb.Assign("id_language", *user.IdLanguage))
		canUpdate = true
	}

	if canUpdate {
		sql, args := sb.Build()
		_, err := db.Exec(sql, args...)
		return err
	}

	return errors.New("no parameters to update User")
}
