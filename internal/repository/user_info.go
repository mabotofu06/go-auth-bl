package repository

import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
	"log"
)

// Coution: 大文字でないと外部パッケージから参照できない
func GetUserInfosByUserId(userId string, db *sql.DB) ([]dto.UserInfo, error) {
	const tableName = "mng_user_info_tbl"
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND user_id = $1", tableName)
	// ユーザIDを元にユーザ情報を取得
	rows, err := db.Query(query, userId)

	if err != nil {
		fmt.Printf("db.Query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var userInfoList []dto.UserInfo

	for rows.Next() {
		var userInfo dto.UserInfo
		if err := rows.Scan(
			&userInfo.UserId,
			&userInfo.UserName,
			&userInfo.Email,
			&userInfo.Phone,
			&userInfo.DeleteFlag,
			&userInfo.CreateDateTime,
			&userInfo.UpdateDateTime,
			&userInfo.DeleteDate,
		); err != nil {
			log.Fatal(err)
		}
		userInfoList = append(userInfoList, userInfo)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err: %v\n", err)
		return nil, err
	}

	return userInfoList, nil
}
