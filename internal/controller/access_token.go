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
	tokenSession, err := session.GetValue[session.TokenInfo](code, req)
	if err != nil {
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

	// キャッシュ初期化 TODO:アクセストークンを取得、設定する処理を共通化して操作しやすくする
	_, e := cache.SetupCache()
	if e != nil {
		fmt.Printf("キャッシュ初期化エラー: %v\n", e)
		middleware.ResError(res, a_err.NewServerErr("内部エラー"))
		return
	}

	// アクセストークンをキャッシュに保存（1時間のTTL）
	ttl := time.Hour * 1
	cache.SetCache(tokenSession.AccessToken, tokenSession, int64(1), ttl)

	fmt.Printf("アクセストークンをキャッシュに保存しました: %s\n", tokenSession.AccessToken)

	body := ResAccessToken{
		AccessToken: tokenSession.AccessToken,
		Expire:      int(time.Now().Add(time.Hour * 1).Unix()),
	}
	ResOk(res, &body)
}
