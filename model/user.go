package model

import (
	"fmt"
	"log"

	"github.com/huandu/go-sqlbuilder"
)

type User struct {
	Id       int    `json:"id"`
	Oauth2Id string `json:"oauth2_id"`
	Enabled  bool   `json:"enabled"`
}

type UserCreate struct {
	Name     string
	Oauth2Id *string
}

func GetUserIdByMatch(idmatch int) (User, error) {
	user := User{}

	sql := fmt.Sprintf(`
		SELECT
			u.id,
			u.oauth2_id,
			u.enabled
		FROM %s AS m
		JOIN user AS u ON m.id_user = u.id
		WHERE m.id = ?
	`, "`match`")

	err := db.QueryRow(sql, idmatch).Scan(&user.Id, &user.Oauth2Id, &user.Enabled)
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
	sb.Select("id", "oauth2_id", "enabled").From("user").Where(sb.Equal("oauth2_id", oauth2Id))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&user.Id, &user.Oauth2Id, &user.Enabled)
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

func CreateUser(oauth2Id string) (int64, error) {
	sb := sqlbuilder.NewInsertBuilder()
	sb.InsertInto("user").Cols("oauth2_id").Values(oauth2Id)
	sql, args := sb.Build()

	response, err := db.Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	id, _ := response.LastInsertId()

	return id, err
}
