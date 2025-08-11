package service

import (
	"fmt"
	"go-auth-bl/internal/dto"
	"go-auth-bl/internal/repository"
	cmn "go-auth-bl/pkg/common"
)

// ユーザIDを元にユーザ認証情報を取得
func IsEnableUserInfo(userId string) error {
	//TODO: このままだと都度DBに接続することになるので、キャッシュを利用するなどの対策を今後考慮
	db := cmn.ConnectDB()
	defer db.Close() // 関数終了時に接続を閉じる

	if exists, err := repository.UserExists(userId, db); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("ユーザ情報が存在しません")
	}

	return nil
}

func InquiryUserInfo(userId string) (*dto.UserInfo, error) {
	db := cmn.ConnectDB()
	defer db.Close()

	userInfo, err := repository.GetUserInfoByUserId(userId, db)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
