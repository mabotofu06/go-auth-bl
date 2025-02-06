package controller

import (
	"fmt"
	"go-auth-bl/def"
	"go-auth-bl/dto"
	apiif "go-auth-bl/dto/if"
	a_err "go-auth-bl/error"
	"go-auth-bl/middleware"
	"go-auth-bl/service"
	"net/http"
	"os"
)

// ログインAPI
func PostLogin(res http.ResponseWriter, req *http.Request) {
	if err := ReqMethodCheck(res, req, POST); err != nil {
		return
	}
	reqBody, err := GetReqBody[apiif.ReqLogin](res, req)
	if err != nil {
		return
	}

	sessionId := req.Header.Get("Session-Id")
	fmt.Printf("sessionId: %s\n", sessionId)
	session, _ := Store.Get(req, "session")

	if val, ok := session.Values[sessionId]; ok {
		authSession := val.(AuthSession)
		fmt.Printf("ClientId: %s, RedirectUri: %s, Scope: %s, State: %s\n",
			authSession.ClientId, authSession.RedirectUri, authSession.Scope, authSession.State)
	} else {
		fmt.Printf("セッションが存在しません\n")
	}

	userAuth, err := getUserAuth(reqBody.UsrId)
	if err != nil {
		middleware.ResError(res, err)
		return
	}
	if err := checkPassword(userAuth, reqBody.Password); err != nil {
		middleware.ResError(res, err)
		return
	}

	// レスポンスを返す
	data := apiif.ResLogin{
		UsrId:   userAuth.UserId,
		Session: "dummy_session",
	}

	ResOk[apiif.ResLogin](res, &data)
}

// サービス層を呼び出してデータを取得
func getUserAuth(uid string) (*dto.UserAuth, *a_err.CustomError) {
	uauth, err := service.GetUserAuthByUserId(uid)
	if err != nil {
		if err == a_err.NotFoundErr {
			return nil, a_err.NewAuthErr("ユーザー名またはパスワードが違います")
		}
		return nil, a_err.NewServerErr(def.ERROR_MESSAGE["E0001"])
	}

	return uauth, nil
}

// パスワードが一致するか確認
func checkPassword(uauth *dto.UserAuth, password string) *a_err.CustomError {
	EncodePassword(password)
	passCheck, err := service.PasswordCheck(uauth, os.Getenv("SALT")+password)
	if err != nil {
		return a_err.NewServerErr(def.ERROR_MESSAGE["E0001"])
	}
	if !passCheck {
		return a_err.NewAuthErr("ユーザー名またはパスワードが違います")
	}

	return nil
}
