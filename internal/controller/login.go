package controller

import (
	"fmt"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/dto"
	apiif "go-auth-bl/internal/dto/if"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	a_err "go-auth-bl/pkg/error"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type TokenSession struct {
	AccessToken string
	RedirectUri string
}

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

	authSession := AuthSession{}

	if val, ok := session.Values[sessionId]; ok {
		authSession = val.(AuthSession)
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

	//認可コードとアクセストークンを発行
	code := uuid.New().String()
	accessToken := uuid.New().String()
	//認可コードでセッションに保存
	session.Values[code] = TokenSession{
		AccessToken: accessToken,
		RedirectUri: authSession.RedirectUri,
	}
	session.Options = &sessions.Options{
		MaxAge: int((30 * time.Second).Seconds()), // 認可コードは30秒の有効期限
	}
	session.Save(req, res)

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
