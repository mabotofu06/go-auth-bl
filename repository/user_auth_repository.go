package repository

import (
	"database/sql"
	"fmt"
	"go-auth-bl/dto"
	"log"
)

// Coution: 大文字でないと外部パッケージから参照できない
func GetUserAuthsByUserId(userId string, db *sql.DB) ([]dto.UserAuth, error) {
	const tableName = "mng_user_auth_tbl"
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND user_id = $1", tableName)
	// ユーザIDを元にユーザ認証情報を取得
	rows, err := db.Query(query, userId)

	if err != nil {
		fmt.Printf("db.Query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var userAuthList []dto.UserAuth

	for rows.Next() {
		var userAuth dto.UserAuth
		if err := rows.Scan(
			&userAuth.UserId,
			&userAuth.Password,
			&userAuth.PasswordHistory1,
			&userAuth.PasswordHistory2,
			&userAuth.PasswordHistory3,
			&userAuth.PasswordFailCnt,
			&userAuth.PasswordLockFlag,
			&userAuth.DeleteFlag,
			&userAuth.CreateDateTime,
			&userAuth.UpdateDateTime,
			&userAuth.DeleteDate,
		); err != nil {
			log.Fatal(err)
		}
		userAuthList = append(userAuthList, userAuth)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err: %v\n", err)
		return nil, err
	}

	return userAuthList, nil
}
