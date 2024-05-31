package model

import (
	"fmt"
	"log"

	"github.com/huandu/go-sqlbuilder"
)

type User struct {
	Id       int    `json:"id"`
	Oauth2Id string `json:"oauth2_id"`
}

type UserCreate struct {
	Name     string
	Oauth2Id *string
}

func GetUsers() ([]User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "oauth2_id").From("user")
	sql, _ := sb.Build()

	rows, err := db.Query(sql)
	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Oauth2Id)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err
}

func GetUser(userId int) (User, error) {
	user := User{}

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "oauth2_id").From("user").Where(sb.Equal("id", userId))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&user.Id, &user.Oauth2Id)

	return user, err
}

func GetUserIdByMatch(idmatch int) (User, error) {
	user := User{}

	sql := fmt.Sprintf(`
		SELECT
			u.id,
			u.oauth2_id
		FROM %s AS m
		JOIN user AS u ON m.id_user = u.id
		WHERE m.id = ?
	`, "`match`")

	err := db.QueryRow(sql, idmatch).Scan(&user.Id, &user.Oauth2Id)

	if err != nil {
		log.Println(err)
	}

	return user, err
}

func GetUserByOauth2Id(oauth2Id string) (User, error) {
	user := User{}

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "oauth2_id").From("user").Where(sb.Equal("oauth2_id", oauth2Id))
	sql, args := sb.Build()

	err := db.QueryRow(sql, args...).Scan(&user.Id, &user.Oauth2Id)

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

func DeleteUser(userId int64) error {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.DeleteFrom("user").Where(sb.Equal("id", userId))
	sql, args := sb.Build()

	_, err := db.Exec(sql, args...)

	return err
}
