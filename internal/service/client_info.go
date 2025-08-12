package service

import (
	"fmt"
	"go-auth-bl/internal/repository"
	cmn "go-auth-bl/pkg/common"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/url"
)

// ユーザIDを元にユーザ認証情報を取得
func IsEnableClient(clientId string, redirectUri string) *ctm_err.CustomError {
	//TODO: このままだと都度DBに接続することになるので、キャッシュを利用するなどの対策を今後考慮
	db, err := cmn.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close() // 関数終了時に接続を閉じる

	// サービス層を呼び出してデータを取得
	clientInfo, err := repository.GetClientInfoByClientId(clientId, db)
	if clientInfo == nil || err != nil {
		return ctm_err.NewDBErr(fmt.Sprintf("クライアント情報取得中にエラー: %v", err))
	}

	logger.Print(logger.DEBUG, "clientInfo: %v", clientInfo)
	logger.Print(logger.DEBUG, "リダイレクト先チェック")

	if err := validateRedirectHost(redirectUri, clientInfo.ClientHost); err != nil {
		return ctm_err.NewRequestErr(fmt.Sprintf("リダイレクト先チェックに失敗: %v", err))
	}

	return nil
}

// リダイレクトURLのホスト名を検証
func validateRedirectHost(redirectUri string, allowedHost string) *ctm_err.CustomError {
	// URLをパースしてホスト名を取得
	parsedURL, err := url.Parse(redirectUri)
	if err != nil {
		logger.Print(logger.ERROR, "無効なリダイレクトURL: %v", err)
		return ctm_err.NewRequestErr("無効なリダイレクトURLです")
	}
	redirectHost := parsedURL.Host
	if redirectHost == "" {
		return ctm_err.NewRequestErr("リダイレクトURLにホスト名が含まれていません")
	}

	// 完全一致でチェック(TODO: 今後サブドメイン含むリダイレクトURIも考慮できるように)
	if redirectHost == allowedHost {
		logger.Print(logger.DEBUG, "ホスト名が一致しました: %s", redirectHost)
		return nil
	}

	logger.Print(logger.WARN, "許可されていないホスト名: %s (許可: %s)", redirectHost, allowedHost)
	return ctm_err.NewRequestErr("許可されていないリダイレクトURLのホスト名です")
}
