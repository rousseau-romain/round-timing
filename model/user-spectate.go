package model

import (
	"log"
)

type UserSpectate struct {
	Id          int    `json:"id"`
	IdUser      int    `json:"id_user"`
	IdUserShare string `json:"id_user_share"`
}

type UserSpectateCreate struct {
	IdUser      int    `json:"id_user"`
	IdUserShare string `json:"id_user_share"`
}

func GetUsersSpectateByIdUser(idUser int) ([]string, error) {

	sql := "SELECT id_user_share FROM user_spectate WHERE id_user = ?"

	rows, err := db.Query(sql, idUser)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	userSpectateIds := []string{}

	if rows.Err() != nil {
		return userSpectateIds, rows.Err()
	}

	for rows.Next() {
		var userSpectateId string
		err := rows.Scan(
			&userSpectateId,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		userSpectateIds = append(userSpectateIds, userSpectateId)
	}

	return userSpectateIds, err
}

func IsUsersSpectateByIdUser(idUser int, idUserShare string) (bool, error) {
	userId := 0

	sql := "SELECT id FROM user_spectate WHERE id_user = ? AND id_user_share = ?"

	err := db.QueryRow(sql, idUser, idUserShare).Scan(&userId)

	if err != nil {
		return false, nil
	}

	return userId != 0, err
}

func CreateUserSpectate(user UserSpectateCreate) (int64, error) {

	sql := "INSERT INTO user_spectate (id_user, id_user_share) VALUES (?, ?)"

	response, err := db.Exec(sql, user.IdUser, user.IdUserShare)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	id, _ := response.LastInsertId()

	return id, err
}

func DeleteUserSpectate(idUser int, idUserShare string) error {
	sql := "DELETE FROM user_spectate WHERE id_user = ? AND id_user_share = ?"

	_, err := db.Exec(sql, idUser, idUserShare)

	return err
}
