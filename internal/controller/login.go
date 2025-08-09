package controller

import (
	"fmt"
	"go-auth-bl/cache"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/dto"
	apiif "go-auth-bl/internal/dto/if"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	"go-auth-bl/internal/session"
	a_err "go-auth-bl/pkg/error"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type ResLogin struct {
	Code        string `json:"code"`
	RedirectUri string `json:"redirectUri"`
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
	// CookieからセッションIDを取得
	cookie, e := req.Cookie("sesid")
	if e != nil {
		fmt.Printf("セッションID取得に失敗しました\n")
		middleware.ResError(res, a_err.NewAuthErr("認可エラー"))
		return
	}
	sessionId := cookie.Value
	fmt.Printf("sessionId: %s\n", sessionId)

	authSession, ok := cache.GetCache[session.PermissionInfo](sessionId, true)

	if !ok {
		fmt.Printf("セッションが存在しません\n")
		middleware.ResError(res, a_err.NewAuthErr("認可エラー"))
		return
	}
	fmt.Printf("authSession: %+v\n", authSession)

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
	tokenSession := session.TokenInfo{
		AccessToken: uuid.New().String(),
		RedirectUri: authSession.RedirectUri,
	}
	if err := cache.SetCache[session.TokenInfo](code, tokenSession, int64(2), 30*time.Minute); err != nil {
		middleware.ResError(res, a_err.NewAuthErr("認可エラー"))
		return
	}

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	data := ResLogin{
		Code:        code,
		RedirectUri: authSession.RedirectUri,
	}

	ResOk[ResLogin](res, &data)

	//res.Header().Set("Location", "/api/v1/redirect"+"?code="+code+"&redirect_uri="+"http://localhost:8080")
	//res.WriteHeader(http.StatusFound) // レスポンスを返す
	//http.Redirect(res, req, "http://localhost:8080"+"?code="+code, http.StatusFound)
	// data := apiif.ResLogin{
	// 	UsrId:   userAuth.UserId,
	// 	Session: "dummy_session",
	// }

	// ResOk[apiif.ResLogin](res, &data)
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
