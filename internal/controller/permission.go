package controller

import (
	"fmt"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/session"
	"net/http"

	"github.com/google/uuid"
)

type ResAuth struct {
	Status string `json:"status"`
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
	cid := queryParams.Get("client_id")       //必須 リクエスト元のクライアントID
	ruri := queryParams.Get("redirect_uri")   //必須　認可サーバはこのURIが登録されている
	scope := queryParams.Get("scope")         //任意　リソースへのアクセス範囲
	state := queryParams.Get("state")         //任意　CSRF対策

	fmt.Printf("rtype=%s\n", rtype)
	fmt.Printf("cid=%s\n", cid)
	fmt.Printf("ruri=%s\n", ruri)
	fmt.Printf("scope=%s\n", scope)
	fmt.Printf("state=%s\n", state)

	// TODO:パラメータチェック
	// if rtype != "code" || cid == "" || ruri == "" {
	// 	middleware.ResError(res, a_err.NewRequestErr("パラメータが不適切です"))
	// 	return
	// }
	//TODO: クライアントIDチェック
	//clientInfo, err := service.GetClientInfo(cid)
	//TODO: リダイレクトURIチェック（DBに登録されたリダイレクト先か判定）

	// セッションID発行
	sessionId := uuid.New().String()
	permission := session.PermissionInfo{
		ClientId:    cid,
		RedirectUri: ruri,
		Scope:       scope,
		State:       state,
	}
	// 認可情報をセッションに保存(有効期限30分)
	if err := session.SetValue[session.PermissionInfo](res, req, sessionId, permission, 30*60); err != nil {
		middleware.ResError(res, err)
		return
	}
	res.Header().Set("sesid", sessionId)

	body := ResAuth{Status: "OK"}
	ResOk(res, &body)
}
