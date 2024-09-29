package service

import (
	"errors"
	"fmt"
	cmn "go-auth-bl/common"
	"go-auth-bl/dto"
	"go-auth-bl/repository"
)

// Coution: 大文字でないと外部パッケージから参照できない
// ユーザIDを元にユーザ認証情報を取得
func GetUserAuthByUserId(userId string) (*dto.UserAuth, error) {
	db := cmn.ConnectDB()
	defer db.Close() // 関数終了時に接続を閉じる

	// サービス層を呼び出してデータを取得
	userAuths, err := repository.GetUserAuthsByUserId(userId, db)
	cmn.ErrorCheck(err)

	fmt.Printf("userAuths: %v\n", userAuths)

	if len(userAuths) == 0 {
		fmt.Println("No user auth data found")
		return nil, errors.New("NotFoundUser")
	}

	return &userAuths[0], nil
}

func AddUserAuth(userId string, password string) (*dto.UserAuth, error) {
	db := cmn.ConnectDB()
	defer db.Close() // 関数終了時に接続を閉じる

	return nil, nil
}
