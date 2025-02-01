package repository

/*
ユーザー認証管理テーブルに対する操作をおこなうリポジトリ
*/

import (
	"database/sql"
	"fmt"
	"go-auth-bl/dto"
	a_err "go-auth-bl/error"
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

/*
認証ユーザ情報取得

@param userId string: ユーザID

@param db *sql.DB: DB接続情報

@throws common_db.NotFoundErr: ユーザ認証情報が見つからない場合

@return *dto.UserAuth: ユーザ認証情報
*/
func GetUserAuthByUserId(userId string, db *sql.DB) (*dto.UserAuth, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND user_id = $1", tableName)
	// ユーザIDを元にユーザ認証情報を取得
	rows, q_err := db.Query(query, userId)
	if q_err != nil {
		fmt.Printf("db.Query: %v\n", q_err)
		return nil, q_err
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
			return nil, err
		}
		userAuthList = append(userAuthList, userAuth)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err: %v\n", err)
		return nil, err
	}
	if len(userAuthList) == 0 {
		fmt.Println("No user auth data found")
		return nil, a_err.NotFoundErr
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
		"  password_lock     = $2,"+
		"  updated_datetime  = $3 "+
		" WHERE "+
		"  user_id = $4", tableName)

	now, date_err := nowDatetime()
	if date_err != nil {
		fmt.Printf("予期せぬエラーが発生しました: %v\n", date_err)
		return date_err
	}
	if _, err := db.Exec(query, failCount, passLock, now, userId); err != nil {
		fmt.Printf("DB更新中にエラーが発生しました: %v\n", err)
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
	if _, err = db.Exec(query, 0, 0, now, userId); err != nil {
		fmt.Printf("DB更新中にエラーが発生しました: %v\n", err)
		return err
	}

	return nil
}
