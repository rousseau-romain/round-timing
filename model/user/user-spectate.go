package user

import "context"

type UserSpectate struct {
	Id          int    `json:"id"`
	IdUser      int    `json:"id_user"`
	IdUserShare string `json:"id_user_share"`
}

type UserSpectateCreate struct {
	IdUser      int    `json:"id_user"`
	IdUserShare string `json:"id_user_share"`
}

func GetUsersSpectateByIdUser(ctx context.Context, idUser int) ([]string, error) {

	sql := "SELECT id_user_share FROM user_spectate WHERE id_user = ?"

	rows, err := db.QueryContext(ctx, sql, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userShareIds []string

	for rows.Next() {
		var userShareId string
		err := rows.Scan(&userShareId)
		if err != nil {
			return nil, err
		}
		userShareIds = append(userShareIds, userShareId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userShareIds, nil
}

func IsUsersSpectateByIdUser(ctx context.Context, idUser int, idUserShare string) (bool, error) {
	userId := 0

	sql := "SELECT id FROM user_spectate WHERE id_user = ? AND id_user_share = ?"

	err := db.QueryRowContext(ctx, sql, idUser, idUserShare).Scan(&userId)

	if err != nil {
		return false, nil
	}

	return userId != 0, err
}

func CreateUserSpectate(ctx context.Context, user UserSpectateCreate) (int64, error) {

	sql := "INSERT INTO user_spectate (id_user, id_user_share) VALUES (?, ?)"

	response, err := db.ExecContext(ctx, sql, user.IdUser, user.IdUserShare)

	if err != nil {
		return 0, err
	}

	id, _ := response.LastInsertId()

	return id, err
}

func DeleteUserSpectate(ctx context.Context, idUser int, idUserShare string) error {
	sql := "DELETE FROM user_spectate WHERE id_user = ? AND id_user_share = ?"

	_, err := db.ExecContext(ctx, sql, idUser, idUserShare)

	return err
}
