package repository

import (
	"database/sql"
	"fmt"
	"go-auth-bl/dto"
	"log"
)

const tableName = "mng_user_auth_tbl "

// Coution: 大文字でないと外部パッケージから参照できない
func GetUserAuthsByUserId(userId string, db *sql.DB) ([]dto.UserAuth, error) {
	// ユーザIDを元にユーザ認証情報を取得

	rows, err := db.Query(
		"SELECT "+
			"* "+
			"FROM "+
			tableName+
			"WHERE "+
			"user_id = $1", userId,
	)

	if err != nil {
		fmt.Printf("db.Query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var userAuths []dto.UserAuth

	for rows.Next() {
		var userAuth dto.UserAuth
		if err := rows.Scan(
			&userAuth.UserId,
			&userAuth.Password,
			&userAuth.SessionToken,
			&userAuth.LastSessinonTime,
			&userAuth.DeleteFlag,
			&userAuth.CreateDateTime,
			&userAuth.UpdateDateTime,
			&userAuth.DeleteDate,
		); err != nil {
			log.Fatal(err)
		}
		userAuths = append(userAuths, userAuth)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err: %v\n", err)
		return nil, err
	}

	return userAuths, nil
}
