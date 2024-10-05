package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-auth-bl/dto"
	"log"
	"time"
)

func nowDatetime() (time.Time, error) {
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, fmt.Errorf("タイムゾーンのロードに失敗しました: %v", err)
	}

	return time.Now().In(location), nil
}

const tableName = "mng_user_auth_tbl"

// 認証ユーザ情報取得
func GetUserAuthByUserId(userId string, db *sql.DB) (*dto.UserAuth, error) {
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
	if len(userAuthList) == 0 {
		fmt.Println("No user auth data found")
		return nil, errors.New("NotFoundUser")
	}

	return &userAuthList[0], nil
}

// パスワード誤り更新
func UpdatePasswordFailNum(userId string, failCount int, db *sql.DB) error {
	PASS_FAIL_CNT := 5
	passLock := 0

	//パスワード誤り回数が指定回数超えていたらロックフラグ
	if failCount >= PASS_FAIL_CNT {
		passLock = 1
	}

	query := fmt.Sprintf("UPDATE %s"+
		" SET"+
		"  password_fail_cnt = $1,"+
		"  password_lock = $2,"+
		"  updated_datetime = $3"+
		" WHERE "+
		"  user_id = $4", tableName)

	now, err := nowDatetime()

	if err != nil {
		fmt.Printf("db.Exec: %v\n", err)
		return err
	}

	_, err = db.Exec(
		query,
		failCount,
		passLock,
		now,
		userId,
	)

	if err != nil {
		fmt.Printf("db.Exec: %v\n", err)
		return err
	}

	return nil
}

// パスワードロック解除
func ResetPasswordLock(userId string, db *sql.DB) error {
	query := fmt.Sprintf("UPDATE %s"+
		" SET"+
		"  password_fail_cnt = $1,"+
		"  password_lock = $2,"+
		"  updated_datetime = $3"+
		" WHERE "+
		"  user_id = $4", tableName)

	now, err := nowDatetime()

	if err != nil {
		fmt.Printf("db.Exec: %v\n", err)
		return err
	}

	_, err = db.Exec(
		query,
		0,
		0,
		now,
		userId,
	)

	if err != nil {
		fmt.Printf("db.Exec: %v\n", err)
		return err
	}

	return nil
}
