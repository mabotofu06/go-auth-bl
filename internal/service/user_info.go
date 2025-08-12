package service

import (
	"go-auth-bl/internal/dto"
	"go-auth-bl/internal/repository"
	cmn "go-auth-bl/pkg/common"
	ctm_err "go-auth-bl/pkg/error"
)

// ユーザIDを元にユーザ認証情報を取得
func IsEnableUserInfo(userId string) *ctm_err.CustomError {
	//TODO: このままだと都度DBに接続することになるので、キャッシュを利用するなどの対策を今後考慮
	db, err := cmn.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close() // 関数終了時に接続を閉じる

	if exists, err := repository.UserExists(userId, db); err != nil {
		return err
	} else if !exists {
		return ctm_err.NewNotFoundErr("ユーザ情報が存在しません")
	}

	return nil
}

func InquiryUserInfo(userId string) (*dto.UserInfo, *ctm_err.CustomError) {
	db, err := cmn.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	userInfo, err := repository.GetUserInfoByUserId(userId, db)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
