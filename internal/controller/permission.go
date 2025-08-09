package controller

import (
	"fmt"
	"go-auth-bl/cache"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	"go-auth-bl/internal/session"
	a_err "go-auth-bl/pkg/error"
	"net/http"
	"time"

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

	if rtype != "code" || cid == "" || ruri == "" {
		middleware.ResError(res, a_err.NewRequestErr("パラメータが不適切です"))
		return
	}
	//クライアントID, リダイレクトURIチェック
	if err := service.IsEnableClient(cid, ruri); err != nil {
		middleware.ResError(res, a_err.NewAuthErr("無効な認証情報です"))
		return
	}

	// セッションID発行
	sessionId := uuid.New().String()
	http.SetCookie(res, &http.Cookie{
		Name:     "sesid",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   30 * 60,
		HttpOnly: true,  //JSからのアクセスを防止（フロント側ではそのまま返却するよう設定）
		Secure:   false, // 開発HTTPなら false に
		SameSite: http.SameSiteLaxMode,
	})
	// 認可情報をセッションに保存(有効期限30分)
	permission := session.PermissionInfo{
		ClientId:    cid,
		RedirectUri: ruri,
		Scope:       scope,
		State:       state,
	}

	if err := cache.SetCache[session.PermissionInfo](sessionId, permission, int64(5), 30*time.Minute); err != nil {
		middleware.ResError(res, a_err.NewAuthErr("セッション保存エラー"))
		return
	}

	body := ResAuth{Status: "OK"}
	ResOk(res, &body)
}
