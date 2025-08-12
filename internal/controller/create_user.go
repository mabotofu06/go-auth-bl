package controller

import (
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
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
	logger.Print(logger.INFO, "API: %sを実行します==========================", def.CONFIG.API["create_user"].Name)

	if err := ReqMethodCheck(res, req, POST); err != nil {
		middleware.ResError(res, err)
		return
	}
	reqBody, err := GetReqBody[reqBodyStruct](res, req)
	if err != nil {
		middleware.ResError(res, err)
		return
	}

	if reqBody.UserId == "" || reqBody.Password == "" {
		middleware.ResError(res, ctm_err.NewRequestErr("userIdまたはpasswordが空です"))
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
