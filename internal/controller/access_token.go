package controller

import (
	"go-auth-bl/internal/cache"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/session"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/http"
	"time"
)

type ResAccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      string `json:"user_id"`
	Expire      int    `json:"expire"`
}

type ReqAccessToken struct {
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
	State       string `json:"state"`
}

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func GetAccessToken(res http.ResponseWriter, req *http.Request) {
	logger.Print(logger.INFO, "API: %sを実行します=========================", def.CONFIG.API["get_token"].Name)
	if err := ReqMethodCheck(res, req, POST); err != nil {
		middleware.ResError(res, err)
		return
	}
	reqBody, err := GetReqBody[ReqAccessToken](res, req)
	if err != nil {
		middleware.ResError(res, err)
		return
	}

	code := reqBody.Code        //必須
	ruri := reqBody.RedirectUri //必須
	//state := reqBody.State      //任意

	// パラメータチェック
	if code == "" || ruri == "" {
		middleware.ResError(res, ctm_err.ParameterErr)
		return
	}

	// ログインAPIで設定したTokenセッション取得
	tokenSession, ok := cache.GetCache[session.CodeInfo](code, true)
	if !ok {
		logger.Print(logger.ERROR, "セッションが存在しません")
		middleware.ResError(res, ctm_err.UnauthorizedErr)
		return
	}

	// リダイレクトURIチェック
	if tokenSession.RedirectUri != ruri {
		logger.Print(logger.ERROR, "リダイレクトURIが不正です")
		middleware.ResError(res, ctm_err.UnauthorizedErr)
		return
	}

	// アクセストークンをキャッシュに保存（1時間の期限）
	ttl := time.Hour * 1
	tokenInfo := session.TokenInfo{
		ClientId: tokenSession.ClientId,
		UserId:   tokenSession.UserId,
		Scope:    tokenSession.Scope,
	}
	cache.SetCache[session.TokenInfo](tokenSession.AccessToken, tokenInfo, int64(1), ttl)

	logger.Print(logger.DEBUG, "アクセストークンをキャッシュに保存しました: %s", tokenSession.AccessToken)

	body := ResAccessToken{
		AccessToken: tokenSession.AccessToken,
		UserId:      tokenSession.UserId,
		Expire:      int(time.Now().Add(time.Hour * 1).Unix()),
	}

	if err := ResOk(res, &body); err != nil {
		middleware.ResError(res, err)
		return
	}
}
