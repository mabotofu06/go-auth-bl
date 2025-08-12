package repository

import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
	"go-auth-bl/pkg/logger"
)

// ユーザIDを元にユーザ情報を取得（プライマリーキーから取得するため1件だけ）
func GetUserInfoByUserId(userId string, db *sql.DB) (*dto.UserInfo, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND user_id = $1", TBL_USER_INFO)
	row := db.QueryRow(query, userId)
	if row == nil {
		return nil, fmt.Errorf("ユーザ情報が見つかりませんでした")
	}

	var userInfo dto.UserInfo
	if err := row.Scan(
		&userInfo.UserId,
		&userInfo.UserName,
		&userInfo.Email,
		&userInfo.Gender,
		&userInfo.Age,
		&userInfo.DeleteFlag,
		&userInfo.CreateDateTime,
		&userInfo.UpdateDateTime,
		&userInfo.DeleteDate,
	); err != nil {
		logger.Print(logger.ERROR, "db.Query: %v", err)
		return nil, err
	}

	return &userInfo, nil
}
