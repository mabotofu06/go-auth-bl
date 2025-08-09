package controller

import (
	"fmt"
	"go-auth-bl/cache"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/session"
	a_err "go-auth-bl/pkg/error"
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
	if ReqMethodCheck(res, req, POST) != nil {
		return
	}
	reqBody, err := GetReqBody[ReqAccessToken](res, req)
	if err != nil {
		return
	}

	code := reqBody.Code        //必須
	ruri := reqBody.RedirectUri //必須
	state := reqBody.State      //任意

	fmt.Printf("code=%s\n", code)
	fmt.Printf("ruri=%s\n", ruri)
	fmt.Printf("state=%s\n", state)

	// パラメータチェック
	if code == "" || ruri == "" {
		middleware.ResError(res, a_err.NewRequestErr("パラメータが不適切です"))
		return
	}

	// ログインAPIで設定したTokenセッション取得
	tokenSession, ok := cache.GetCache[session.CodeInfo](code, true)
	if !ok {
		fmt.Printf("セッションが存在しません\n")
		middleware.ResError(res, a_err.NewAuthErr("認可エラー"))
		return
	}

	// リダイレクトURIチェック
	if tokenSession.RedirectUri != ruri {
		fmt.Printf("リダイレクトURIが不正です\n")
		middleware.ResError(res, a_err.NewAuthErr("認可エラー"))
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

	fmt.Printf("アクセストークンをキャッシュに保存しました: %s\n", tokenSession.AccessToken)

	body := ResAccessToken{
		AccessToken: tokenSession.AccessToken,
		UserId:      tokenSession.UserId,
		Expire:      int(time.Now().Add(time.Hour * 1).Unix()),
	}
	ResOk(res, &body)
}
