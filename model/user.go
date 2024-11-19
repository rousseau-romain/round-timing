package model

import (
	"errors"
	"fmt"
	"log"

	"github.com/huandu/go-sqlbuilder"
)

type User struct {
	Id       int    `json:"id"`
	Oauth2Id string `json:"oauth2_id"`
	Enabled  bool   `json:"enabled"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}
type UserUpdate struct {
	Enabled *bool
}

type UserCreate struct {
	Oauth2Id string
	Email    string
}

func GetUsers() ([]User, error) {
	sql := `
		SELECT 
			id,
			email,
			enabled,
			is_admin
		FROM user
	`

	rows, err := db.Query(sql)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Enabled,
			&user.IsAdmin,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}

func GetUserIdByMatch(idmatch int) (User, error) {
	user := User{}

	sql := fmt.Sprintf(`
		SELECT
			u.id,
			u.oauth2_id,
			u.enabled,
			u.email,
			u.is_admin
		FROM %s AS m
		JOIN user AS u ON m.id_user = u.id
		WHERE m.id = ?
	`, "`match`")

	err := db.QueryRow(sql, idmatch).Scan(&user.Id, &user.Oauth2Id, &user.Enabled, &user.Email, &user.IsAdmin)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user, nil
	}
	if err != nil {
		log.Println(err)
	}

	return user, err
}

func GetUserByOauth2Id(oauth2Id string) (User, error) {
	user := User{}

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "oauth2_id", "enabled", "email", "is_admin").From("user").Where(sb.Equal("oauth2_id", oauth2Id))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&user.Id, &user.Oauth2Id, &user.Enabled, &user.Email, &user.IsAdmin)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return user, nil
	}
	return user, err
}

func UserExists(oauth2Id string) (bool, error) {
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
		INSERT INTO user (oauth2_id, email)
		VALUES ("?", "?")
	`

	response, err := db.Exec(sql, user.Oauth2Id, user.Email)

	if err != nil {
		log.Println(err)
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

	if canUpdate {
		sql, args := sb.Build()
		_, err := db.Exec(sql, args...)
		return err
	}

	return errors.New("no parameters to update User")
}
