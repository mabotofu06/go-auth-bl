package controller

import (
	"fmt"
	"go-auth-bl/cache"
	"go-auth-bl/internal/middleware"
	a_err "go-auth-bl/pkg/error"
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
	if ReqMethodCheck(res, req, GET) != nil {
		return
	}
	queryParams := req.URL.Query()
	tkn := queryParams.Get("token") //必須 固定値:"token"

	fmt.Printf("token=%s\n", tkn)

	if tkn == "" {
		middleware.ResError(res, a_err.NewRequestErr("パラメータが不適切です"))
		return
	}

	tokenInfo, err := cache.GetTokenFromCache(tkn)
	if err != nil {
		middleware.ResError(res, a_err.NewAuthErr("無効なトークンです"))
		return
	}

	body := ResToken{Token: tokenInfo.AccessToken}
	ResOk(res, &body)
}
