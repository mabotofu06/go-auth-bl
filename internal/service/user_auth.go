package service

import (
	"fmt"
	"go-auth-bl/internal/dto"
	"go-auth-bl/internal/repository"
	cmn "go-auth-bl/pkg/common"

	"golang.org/x/crypto/bcrypt"
)

// ユーザIDを元にユーザ認証情報を取得
func GetUserAuthByUserId(userId string) (*dto.UserAuth, error) {
	db := cmn.ConnectDB()
	defer db.Close() // 関数終了時に接続を閉じる

	// サービス層を呼び出してデータを取得
	userAuth, err := repository.GetUserAuthByUserId(userId, db)
	if err != nil {
		return nil, err
	}
	fmt.Printf("userAuth: %v\n", userAuth)
	return userAuth, nil
}

// パスワードが一致するか確認
func PasswordCheck(userAuth *dto.UserAuth, password string) (bool, error) {
	db := cmn.ConnectDB()
	defer db.Close() // 関数終了時に接続を閉じる

	if userAuth.PasswordLockFlag != 0 {
		fmt.Println("パスワードがロックされています")
		return false, nil
	}

	failCnt := userAuth.PasswordFailCnt

	if !ComparePassword(userAuth.Password, password) {
		failCnt++
		//テーブルに対してパスワード失敗回数を加算
		fmt.Println("パスワードが一致しませんでした")
		fmt.Printf("failCnt: %d\n", failCnt)

		//パスワードロック込みでDBを更新
		err := repository.UpdatePasswordFailNum(userAuth.UserId, failCnt, db)

		return false, err
	}

	fmt.Println("Password matches")
	//テーブルに対してパスワード失敗回数をリセット
	err := repository.ResetPasswordLock(userAuth.UserId, db)
	return true, err
}

// パスワードを比較する関数
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func AddUserAuth(userId string, password string) (*dto.UserAuth, error) {
	db := cmn.ConnectDB()
	defer db.Close() // 関数終了時に接続を閉じる

	return nil, nil
}
