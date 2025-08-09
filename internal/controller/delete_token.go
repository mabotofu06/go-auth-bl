package controller

import (
	"fmt"
	"go-auth-bl/cache"
	"go-auth-bl/internal/middleware"
	a_err "go-auth-bl/pkg/error"
	"net/http"
)

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func DeleteToken(res http.ResponseWriter, req *http.Request) {
	if ReqMethodCheck(res, req, DELETE) != nil {
		return
	}

	reqBody, err := GetReqBody[struct {
		Token string `json:"token"`
	}](res, req)
	if err != nil {
		return
	}
	tkn := reqBody.Token

	fmt.Printf("token=%s\n", tkn)

	if tkn == "" {
		middleware.ResError(res, a_err.NewRequestErr("パラメータが不適切です"))
		return
	}

	if err := cache.DeleteCache(tkn); err != nil {
		middleware.ResError(res, a_err.NewServerErr("トークン削除中エラーが発生しました。"))
		return
	}

	body := struct {
		Message string `json:"message"`
	}{
		Message: "トークンが正常に削除されました。",
	}
	ResOk(res, &body)
}
