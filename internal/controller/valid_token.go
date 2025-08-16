package controller

import (
	"go-auth-bl/internal/cache"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/session"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/http"
)

type ResToken struct {
	Token string `json:"token"`
}

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func GetValidToken(res http.ResponseWriter, req *http.Request) {
	logger.Print(logger.INFO, "API: %sを実行します===================================", def.CONFIG.API["check_token"].Name)

	if err := ReqMethodCheck(res, req, GET); err != nil {
		middleware.ResError(res, err)
		return
	}
	queryParams := req.URL.Query()
	tkn := queryParams.Get("token") //必須

	if tkn == "" {
		middleware.ResError(res, ctm_err.ParameterErr)
		return
	}

	logger.Print(logger.DEBUG, "クエリパラメータ: %v", queryParams)

	if _, ok := cache.GetCache[session.TokenInfo](tkn, false); !ok {
		middleware.ResError(res, ctm_err.NewAuthErr("無効なトークンです"))
		return
	}

	body := ResToken{Token: tkn}
	if err := ResOk[ResToken](res, &body); err != nil {
		middleware.ResError(res, err)
		return
	}
}
