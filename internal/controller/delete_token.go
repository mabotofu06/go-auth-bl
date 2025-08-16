package controller

import (
	"go-auth-bl/internal/cache"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/middleware"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/http"
)

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func DeleteToken(res http.ResponseWriter, req *http.Request) {
	logger.Print(logger.INFO, "API: %sを実行します===========================", def.CONFIG.API["delete_token"].Name)

	if err := ReqMethodCheck(res, req, DELETE); err != nil {
		middleware.ResError(res, err)
		return
	}

	reqBody, err := GetReqBody[struct {
		Token string `json:"token"`
	}](res, req)
	if err != nil {
		middleware.ResError(res, err)
		return
	}
	tkn := reqBody.Token

	if tkn == "" {
		middleware.ResError(res, ctm_err.ParameterErr)
		return
	}

	if err := cache.DeleteCache(tkn); err != nil {
		middleware.ResError(res, ctm_err.UnexpectedServerErr)
		return
	}

	body := struct {
		Message string `json:"message"`
	}{
		Message: "トークンが正常に削除されました。",
	}
	if err := ResOk(res, &body); err != nil {
		middleware.ResError(res, err)
		return
	}
}
