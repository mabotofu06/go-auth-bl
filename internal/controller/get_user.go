package controller

import (
	"go-auth-bl/cache"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	"go-auth-bl/internal/session"
	a_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/http"
)

// ユーザ情報取得API
func GetUserInfo(res http.ResponseWriter, req *http.Request) {
	logger.Print(logger.INFO, "API: %sを実行します============================", def.CONFIG.API["get_user"].Name)

	if err := ReqMethodCheck(res, req, GET); err != nil {
		middleware.ResError(res, err)
		return
	}

	sessionId := req.Header.Get("session")
	tokenInfo, ok := cache.GetCache[session.TokenInfo](sessionId, false)
	if !ok {
		middleware.ResError(res, a_err.NewRequestErr("セッション情報が取得できません"))
		return
	}
	//クエリのユーザIDとセッションのユーザIDが異ならないかチェック
	query := req.URL.Query()
	if query.Get("user_id") == "" || query.Get("user_id") != tokenInfo.UserId {
		middleware.ResError(res, a_err.NewRequestErr("権限がありません"))
		return
	}

	//ユーザ情報を取得
	userInfo, err := service.InquiryUserInfo(tokenInfo.UserId)
	if err != nil {
		middleware.ResError(res, a_err.NewServerErr("ユーザ情報の取得に失敗しました"))
		return
	}

	resBody := struct {
		UserId   string  `json:"userId"`
		UserName string  `json:"userName"`
		Email    *string `json:"email"`
	}{
		UserId:   userInfo.UserId,
		UserName: userInfo.UserName,
		Email:    userInfo.Email,
	}

	if err := ResOk(res, &resBody); err != nil {
		middleware.ResError(res, err)
		return
	}
}
