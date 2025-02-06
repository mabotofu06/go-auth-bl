package controller

import (
	"fmt"
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

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func GetAccessToken(res http.ResponseWriter, req *http.Request) {
	if ReqMethodCheck(res, req, GET) != nil {
		return
	}
	queryParams := req.URL.Query()
	code := queryParams.Get("code")         //必須
	ruri := queryParams.Get("redirect_uri") //必須　認可サーバはこのURIが登録されている
	state := queryParams.Get("state")       //任意　CSRF対策

	fmt.Printf("code=%s\n", code)
	fmt.Printf("ruri=%s\n", ruri)
	fmt.Printf("state=%s\n", state)

	// パラメータチェック
	// if code != "" || ruri == "" {
	// 	middleware.ResError(res, a_err.NewRequestErr("パラメータが不適切です"))
	// 	return
	// }

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

	// TODO:アクセストークンをDBに保存

	body := ResAccessToken{
		AccessToken: tokenSession.AccessToken,
		Expire:      int(time.Now().Add(time.Hour * 1).Unix()),
	}
	ResOk(res, &body)
}
