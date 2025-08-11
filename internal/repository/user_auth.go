package repository

/*
ユーザー認証管理テーブルに対する操作をおこなうリポジトリ
*/

import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
	"time"
)

func nowDatetime() (time.Time, error) {
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, fmt.Errorf("タイムゾーンのロードに失敗しました: %v", err)
	}
	return time.Now().In(location), nil
}

/*
認証ユーザ情報取得

@param userId string: ユーザID

@param db *sql.DB: DB接続情報

@throws common_db.NotFoundErr: ユーザ認証情報が見つからない場合

@return *dto.UserAuth: ユーザ認証情報
*/
func GetUserAuthByUserId(userId string, db *sql.DB) (*dto.UserAuth, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND user_id = $1", TBL_USER_AUTH)
	// ユーザIDを元にユーザ認証情報を取得
	row := db.QueryRow(query, userId)
	if row == nil {
		return nil, fmt.Errorf("ユーザ情報が見つかりませんでした")
	}

	var userAuth dto.UserAuth
	if err := row.Scan(
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
		fmt.Printf("db.Query: %v\n", err)
		return nil, err
	}

	return &userAuth, nil
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
		"  user_id = $4", TBL_USER_AUTH)

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
		"  user_id = $4", TBL_USER_AUTH)

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

// ユーザ存在確認
func UserExists(userId string, db *sql.DB) (bool, error) {
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM %s WHERE delete_flag = 0 AND user_id = $1)",
		TBL_USER_AUTH,
	)
	var exists bool
	if err := db.QueryRow(query, userId).Scan(&exists); err != nil {
		return false, fmt.Errorf("ユーザ存在確認中にエラー: %w", err)
	}
	return exists, nil
}

func InsertUserAuth(userId string, hashedPassword string, db *sql.DB) (*dto.UserAuth, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, password) VALUES ($1, $2) RETURNING *", TBL_USER_AUTH)
	row := db.QueryRow(query, userId, hashedPassword)

	var userAuth dto.UserAuth
	if err := row.Scan(
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
		return nil, fmt.Errorf("DB挿入中にエラーが発生しました: %v", err)
	}

	return &userAuth, nil
}

// 認証とユーザー情報を1トランザクションで作成
func CreateUserWithInfo(db *sql.DB, userId string, userName string, hashedPassword string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		// 異常時はロールバック（コミット成功時は no-op）
		_ = tx.Rollback()
	}()

	// 認証テーブルへ登録
	queryAuth := fmt.Sprintf("INSERT INTO %s (user_id, password) VALUES ($1, $2)", TBL_USER_AUTH)
	if _, err := tx.Exec(queryAuth, userId, hashedPassword); err != nil {
		return fmt.Errorf("insert auth: %w", err)
	}

	// ユーザー情報テーブルへ登録（email/phone は未入力なら空で）
	queryInfo := fmt.Sprintf("INSERT INTO %s (user_id, user_name) VALUES ($1, $2)", TBL_USER_INFO)
	if _, err := tx.Exec(queryInfo, userId, userName); err != nil {
		return fmt.Errorf("insert user_info: %w", err)
	}

	// すべて成功したらコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
