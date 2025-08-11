package controller

import (
	"fmt"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	a_err "go-auth-bl/pkg/error"
	"net/http"
)

type reqBodyStruct struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type resBodyStruct struct {
	UserId string `json:"userId"`
}

// ユーザ登録API
func PostCreateUser(res http.ResponseWriter, req *http.Request) {
	if err := ReqMethodCheck(res, req, POST); err != nil {
		middleware.ResError(res, err)
		return
	}

	reqBody, err := GetReqBody[reqBodyStruct](res, req)
	if err != nil {
		middleware.ResError(res, err)
		return
	}

	fmt.Printf("reqBody: %v\n", reqBody)
	if reqBody.UserId == "" || reqBody.Password == "" {
		middleware.ResError(res, a_err.NewRequestErr("userIdまたはpasswordが空です"))
		return
	}

	//ログインユーザとユーザ情報を登録
	userId, err := service.CreateNewLoginUser(reqBody.UserId, reqBody.UserName, reqBody.Password)
	if err != nil {
		middleware.ResError(res, err)
		return
	}

	data := resBodyStruct{
		UserId: *userId,
	}

	if err := ResOk[resBodyStruct](res, &data); err != nil {
		middleware.ResError(res, err)
		return
	}
}
