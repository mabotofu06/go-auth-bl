package service

import (
	"go-auth-bl/internal/dto"
	"go-auth-bl/internal/repository"
	cmn "go-auth-bl/pkg/common"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// ユーザIDを元にユーザ認証情報を取得
func GetUserAuthByUserId(userId string) (*dto.UserAuth, *ctm_err.CustomError) {
	db, err := cmn.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close() // 関数終了時に接続を閉じる

	// サービス層を呼び出してデータを取得
	userAuth, err := repository.GetUserAuthByUserId(userId, db)
	if err != nil {
		return nil, err
	}
	logger.Print(logger.INFO, "userAuth: %v", userAuth)
	return userAuth, nil
}

// パスワードが一致するか確認
func PasswordCheck(userAuth *dto.UserAuth, password string) (bool, *ctm_err.CustomError) {
	db, err := cmn.ConnectDB()
	if err != nil {
		return false, err
	}
	defer db.Close() // 関数終了時に接続を閉じる

	if userAuth.PasswordLockFlag != 0 {
		logger.Print(logger.WARN, "パスワードがロックされています")
		return false, nil
	}

	failCnt := userAuth.PasswordFailCnt

	if !ComparePassword(userAuth.Password, password) {
		failCnt++
		//テーブルに対してパスワード失敗回数を加算
		logger.Print(logger.WARN, "パスワードが一致しませんでした")
		logger.Print(logger.INFO, "failCnt: %d", failCnt)

		//パスワードロック込みでDBを更新
		err := repository.UpdatePasswordFailNum(userAuth.UserId, failCnt, db)

		return false, err
	}

	logger.Print(logger.INFO, "Password matches")
	//テーブルに対してパスワード失敗回数をリセット
	err = repository.ResetPasswordLock(userAuth.UserId, db)
	return true, err
}

// パスワードを比較する関数
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// 新規ログインユーザを作成
func CreateNewLoginUser(userId string, userName string, password string) (*string, *ctm_err.CustomError) {
	db, err := cmn.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close() // 関数終了時に接続を閉じる

	//登録済かチェック
	exists, err := repository.UserExists(userId, db)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ctm_err.NewRequestErr("ユーザIDがすでに存在します")
	}

	hashedPassword, e := bcrypt.GenerateFromPassword([]byte(os.Getenv("SALT")+password), bcrypt.DefaultCost)
	if e != nil {
		return nil, ctm_err.UnexpectedServerErr
	}

	// 新規ユーザを作成
	if err := repository.CreateUserWithInfo(db, userId, userName, string(hashedPassword)); err != nil {
		return nil, err
	}

	return &userId, nil
}
