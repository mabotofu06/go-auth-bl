package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type AuthSession struct {
	clientId    string
	redirectUri string
	scope       string
	state       string
}

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func GetPermission(res http.ResponseWriter, req *http.Request) {
	if ReqMethodCheck(res, req, GET) != nil {
		return
	}
	queryParams := req.URL.Query()
	rtype := queryParams.Get("response_type") //必須 固定値:"code"
	cid := queryParams.Get("crient_id")       //必須 リクエスト元のクライアントID
	ruri := queryParams.Get("redirect_uri")   //必須　認可サーバはこのURIが登録されている
	scope := queryParams.Get("scope")         //任意　リソースへのアクセス範囲
	state := queryParams.Get("state")         //任意　CSRF対策

	fmt.Printf("rtype=%s\n", rtype)
	fmt.Printf("cid=%s\n", cid)
	fmt.Printf("ruri=%s\n", ruri)
	fmt.Printf("scope=%s\n", scope)
	fmt.Printf("state=%s\n", state)

	// パラメータチェック
	// if rtype != "code" || cid == "" || ruri == "" {
	// 	middleware.ResError(res, a_err.NewRequestErr("パラメータが不適切です"))
	// 	return
	// }
	//TODO: クライアントIDチェック
	//clientInfo, err := service.GetClientInfo(cid)

	//TODO: リダイレクトURIチェック

	// スコープチェック
	// ステートチェック

	//redirect_uri := queryParams.Get("redirect_uri")
	//fmt.Printf("認可コード発行処理を開始します code=%s\n", code)

	// Basic認証
	req.Header.Get("Authorization")

	// セッションID発行
	sessionId := uuid.New().String()
	session, _ := Store.Get(req, "session")
	// 認可情報をセッションに保存
	session.Values[sessionId] = AuthSession{
		clientId:    cid,
		redirectUri: ruri,
		scope:       scope,
		state:       state,
	}
	res.Header().Set("sesid", sessionId)
	// UIを返却（ログイン画面）
	http.FileServer(http.Dir("./build")).ServeHTTP(res, req)
}
