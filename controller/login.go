package controller

import (
	"fmt"
	apiif "go-auth-bl/dto/if"
	a_err "go-auth-bl/error"
	"go-auth-bl/service"
	"net/http"
	"os"
)

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func PostLogin(res http.ResponseWriter, req *http.Request) {
	ReqMethodCheck(req, POST)

	reqBody := GetReqBody[apiif.ReqLogin](req)

	// サービス層を呼び出してデータを取得
	userAuth, err := service.GetUserAuthByUserId(reqBody.UsrId)
	if err != nil {
		if err == a_err.NotFoundErr {
			a_err.Throw(a_err.NewAuthErr("ユーザー名またはパスワードが違います"))
		}
		a_err.Throw(a_err.NewServerErr("予期せぬエラーが発生しました"))
	}

	// パスワードが一致するか確認
	EncodePassword(reqBody.Password)
	passCheck, err := service.PasswordCheck(userAuth, os.Getenv("SALT")+reqBody.Password)
	if err != nil {
		a_err.Throw(a_err.NewServerErr("予期せぬエラーが発生しました"))
	}
	if !passCheck {
		a_err.Throw(a_err.NewAuthErr("ユーザー名またはパスワードが違います"))
	}

	fmt.Printf("ユーザ情報: %v\n", userAuth)
	// レスポンスを返す
	data := apiif.ResLogin{
		UsrId:   userAuth.UserId,
		Session: "dummy_session",
	}

	ResOk[apiif.ResLogin](res, &data)
}
