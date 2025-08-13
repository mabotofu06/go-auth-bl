package repository

/*
ユーザー認証管理テーブルに対する操作をおこなうリポジトリ
*/

import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
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
func GetUserAuthByUserId(userId string, db *sql.DB) (*dto.UserAuth, *ctm_err.CustomError) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND user_id = $1", TBL_USER_AUTH)
	// ユーザIDを元にユーザ認証情報を取得
	row := db.QueryRow(query, userId)
	if row == nil {
		return nil, ctm_err.NewNotFoundErr("ユーザ情報が見つかりませんでした")
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
		logger.Print(logger.ERROR, "db.Query: %v", err)
		return nil, ctm_err.UnexpectedDBErr
	}

	return &userAuth, nil
}

// パスワード誤り更新
func UpdatePasswordFailNum(userId string, failCount int, db *sql.DB) *ctm_err.CustomError {
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
		logger.Print(logger.ERROR, "予期せぬエラーが発生しました: %v", date_err)
		return ctm_err.UnexpectedServerErr
	}
	if _, err := db.Exec(query, failCount, passLock, now, userId); err != nil {
		logger.Print(logger.ERROR, "DB更新中にエラーが発生しました: %v", err)
		return ctm_err.UnexpectedDBErr
	}

	return nil
}

// パスワードロック解除
func ResetPasswordLock(userId string, db *sql.DB) *ctm_err.CustomError {
	query := fmt.Sprintf("UPDATE %s"+
		" SET"+
		"  password_fail_cnt = $1,"+
		"  password_lock = $2,"+
		"  updated_datetime = $3"+
		" WHERE "+
		"  user_id = $4", TBL_USER_AUTH)

	now, err := nowDatetime()
	if err != nil {
		logger.Print(logger.ERROR, "db.Exec: %v", err)
		return ctm_err.UnexpectedServerErr
	}
	if _, err = db.Exec(query, 0, 0, now, userId); err != nil {
		logger.Print(logger.ERROR, "DB更新中にエラーが発生しました: %v", err)
		return ctm_err.UnexpectedDBErr
	}

	return nil
}

// ユーザ存在確認
func UserExists(userId string, db *sql.DB) (bool, *ctm_err.CustomError) {
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM %s WHERE delete_flag = 0 AND user_id = $1)",
		TBL_USER_AUTH,
	)
	var exists bool
	if err := db.QueryRow(query, userId).Scan(&exists); err != nil {
		return false, ctm_err.UnexpectedDBErr
	}
	return exists, nil
}

func InsertUserAuth(userId string, hashedPassword string, db *sql.DB) (*dto.UserAuth, *ctm_err.CustomError) {
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
		return nil, ctm_err.UnexpectedDBErr
	}

	return &userAuth, nil
}

// 認証とユーザー情報を1トランザクションで作成
func CreateUserWithInfo(db *sql.DB, userId string, userName string, hashedPassword string) *ctm_err.CustomError {
	tx, err := db.Begin()
	if err != nil {
		return ctm_err.NewDBErr("トランザクション開始中にエラーが発生しました")
	}
	defer func() {
		// 異常時はロールバック（コミット成功時は no-op）
		if err == nil {
			logger.Print(logger.INFO, "トランザクションが正常に完了しました")
			return
		}
		logger.Print(logger.ERROR, "異常が発生したためトランザクションをロールバックします")
		_ = tx.Rollback()
	}()

	// 認証テーブルへ登録
	queryAuth := fmt.Sprintf("INSERT INTO %s (user_id, password, admin) VALUES ($1, $2, $3)", TBL_USER_AUTH)
	//権限は後々登録できるようにとりあえず0(一般ユーザ)
	if _, err := tx.Exec(queryAuth, userId, hashedPassword, 0); err != nil {
		return ctm_err.UnexpectedDBErr
	}

	// ユーザー情報テーブルへ登録（email/phone は未入力なら空で）
	queryInfo := fmt.Sprintf("INSERT INTO %s (user_id, user_name) VALUES ($1, $2)", TBL_USER_INFO)
	if _, err := tx.Exec(queryInfo, userId, userName); err != nil {
		return ctm_err.UnexpectedDBErr
	}

	// すべて成功したらコミット
	if err := tx.Commit(); err != nil {
		return ctm_err.UnexpectedDBErr
	}
	return nil
}
