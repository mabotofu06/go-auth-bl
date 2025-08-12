package controller

import (
	"go-auth-bl/internal/cache"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	"go-auth-bl/internal/session"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
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
	logger.Print(logger.INFO, "API: %sを実行します====================================", def.CONFIG.API["permission"].Name)

	if err := ReqMethodCheck(res, req, GET); err != nil {
		middleware.ResError(res, err)
		return
	}
	queryParams := req.URL.Query()
	logger.Print(logger.DEBUG, "クエリパラメータ: %v", queryParams)

	rtype := queryParams.Get("response_type") //必須 固定値:"code"
	cid := queryParams.Get("client_id")       //必須 リクエスト元のクライアントID
	ruri := queryParams.Get("redirect_uri")   //必須　認可サーバはこのURIが登録されている
	scope := queryParams.Get("scope")         //任意　リソースへのアクセス範囲
	state := queryParams.Get("state")         //任意　CSRF対策

	//パラメータチェック
	if rtype != "code" || cid == "" || ruri == "" {
		middleware.ResError(res, ctm_err.ParameterErr)
		return
	}
	//クライアントID, リダイレクトURIチェック
	if err := service.IsEnableClient(cid, ruri); err != nil {
		middleware.ResError(res, ctm_err.NewAuthErr("無効な認証情報です"))
		return
	}

	sessionId := uuid.New().String()
	if c, _ := req.Cookie("sesid"); c != nil && c.Value != "" {
		// 既にセッションIDが存在する場合、既存セッションを再利用
		logger.Print(logger.DEBUG, "既存セッションを再利用します: %s\n", c.Value)
		sessionId = c.Value
	}

	// 新規登録 or 有効期限を延長
	session.SetSessionId(res, sessionId)

	permission := session.PermissionInfo{
		ClientId:    cid,
		RedirectUri: ruri,
		Scope:       scope,
		State:       state,
	}

	// 認可情報をセッションに保存(有効期限30分)
	//ログイン情報入力して送信まで30分有効期限を設ける
	if err := cache.SetCache[session.PermissionInfo](sessionId, permission, int64(5), 30*time.Minute); err != nil {
		middleware.ResError(res, ctm_err.NewAuthErr("セッションエラー"))
		return
	}

	body := ResAuth{Status: "OK"}
	if err := ResOk(res, &body); err != nil {
		middleware.ResError(res, err)
		return
	}
}
