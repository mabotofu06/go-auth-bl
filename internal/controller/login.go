package controller

import (
	"go-auth-bl/internal/cache"
	"go-auth-bl/internal/def"
	"go-auth-bl/internal/dto"
	apiif "go-auth-bl/internal/dto/if"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/service"
	"go-auth-bl/internal/session"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
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
	logger.Print(logger.INFO, "API: %sを実行します==============================", def.CONFIG.API["login"].Name)

	if err := ReqMethodCheck(res, req, POST); err != nil {
		middleware.ResError(res, err)
		return
	}
	reqBody, err := GetReqBody[apiif.ReqLogin](res, req)
	if err != nil {
		middleware.ResError(res, err)
		return
	}
	// CookieからセッションIDを取得
	cookie, e := req.Cookie("sesid")
	if e != nil {
		logger.Print(logger.ERROR, "セッションIDが取得できません: %v", e)
		middleware.ResError(res, ctm_err.UnauthorizedErr)
		return
	}
	sessionId := cookie.Value
	logger.Print(logger.DEBUG, "sessionId: %s", sessionId)

	//1回でもログインミスしたらセッションが無くなって認可エラーとなるため維持する（5回間違えればロックかかるため総攻撃では突破できない想定）
	permissionInfo, ok := cache.GetCache[session.PermissionInfo](sessionId, false)
	if !ok {
		logger.Print(logger.ERROR, "セッションが存在しません")
		middleware.ResError(res, ctm_err.UnauthorizedErr)
		return
	}
	logger.Print(logger.DEBUG, "permissionInfo: %+v", permissionInfo)

	userAuth, err := getUserAuth(reqBody.UsrId)
	if err != nil {
		middleware.ResError(res, err)
		return
	}
	if err := checkPassword(userAuth, reqBody.Password); err != nil {
		middleware.ResError(res, err)
		return
	}

	//認可コード, スコープ, アクセストークンを発行
	code := uuid.New().String()
	tokenSession := session.CodeInfo{
		AccessToken: uuid.New().String(),
		ClientId:    permissionInfo.ClientId,
		UserId:      userAuth.UserId,
		Scope:       permissionInfo.Scope,
		RedirectUri: permissionInfo.RedirectUri,
	}
	//アクセストークン要求まで30秒以内に完了される想定でセッションに保存
	if err := cache.SetCache[session.CodeInfo](code, tokenSession, int64(2), 30*time.Second); err != nil {
		middleware.ResError(res, ctm_err.UnauthorizedErr)
		return
	}

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	data := ResLogin{
		Code:        code,
		RedirectUri: permissionInfo.RedirectUri,
	}

	if err := ResOk[ResLogin](res, &data); err != nil {
		middleware.ResError(res, err)
		return
	}
}

// サービス層を呼び出してデータを取得
func getUserAuth(uid string) (*dto.UserAuth, *ctm_err.CustomError) {
	uauth, err := service.GetUserAuthByUserId(uid)
	if err != nil {
		if err == ctm_err.NotFoundErr {
			return nil, ctm_err.NewAuthErr("ユーザー名またはパスワードが違います")
		}
		return nil, ctm_err.NewServerErr(def.ERROR_MESSAGE["E0001"])
	}

	return uauth, nil
}

// パスワードが一致するか確認
func checkPassword(uauth *dto.UserAuth, password string) *ctm_err.CustomError {
	EncodePassword(password)
	passCheck, err := service.PasswordCheck(uauth, os.Getenv("SALT")+password)
	if err != nil {
		return ctm_err.NewServerErr(def.ERROR_MESSAGE["E0001"])
	}
	if !passCheck {
		return ctm_err.NewAuthErr("ユーザー名またはパスワードが違います")
	}

	return nil
}
