package controller

import (
	"encoding/json"
	"fmt"
	"go-auth-bl/internal/cache"
	apiif "go-auth-bl/internal/dto/if"
	"go-auth-bl/internal/session"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// リクエストメソッドが不適切な場合はエラーを返す
func ReqMethodCheck(res http.ResponseWriter, req *http.Request, method string) *ctm_err.CustomError {
	if req.Method != method {
		err := ctm_err.NewRequestErr("リクエストメソッドが不適切です")
		logger.Print(logger.ERROR, "リクエストメソッドが不適切です: %s", req.Method)
		return err
	}
	return nil
}

// POSTリクエストのリクエストボディを取得（リクエストボディが不適切な場合はエラーを返す）
func GetReqBody[T any](res http.ResponseWriter, req *http.Request) (*T, *ctm_err.CustomError) {
	if req.Body == nil {
		err := ctm_err.NewRequestErr("リクエストボディが空です")
		logger.Print(logger.ERROR, "リクエストボディが空です")
		return nil, err
	}

	defer req.Body.Close()
	var request T

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		err := ctm_err.NewRequestErr("リクエストボディが不適切です")
		logger.Print(logger.ERROR, "リクエストボディエンコード中にエラーが発生しました: %v", err)
		return nil, err
	}

	logger.Print(logger.DEBUG, "リクエストボディ内容: %+v", request)
	return &request, nil
}

// APIの正常終了時のレスポンスを返す
func ResOk[T any](res http.ResponseWriter, data *T) *ctm_err.CustomError {
	logger.Print(logger.DEBUG, "レスポンスデータ: %+v", data)

	res.Header().Set("Content-Type", "application/json")
	resBody := apiif.Response[T]{
		Status: http.StatusOK,
		Code:   "I0001",
		Type:   "正常",
		Msg:    "通信が正常終了しました",
		Data:   data,
	}

	json, err := json.Marshal(resBody)
	if err != nil {
		return ctm_err.NewServerErr("予期せぬエラーが発生しました")
	}
	res.WriteHeader(http.StatusOK)
	res.Write(json)
	return nil
}

// bcryptを使ってパスワードをハッシュ化
// エンコードされた文字列の長さは60文字
func EncodePassword(password string) (string, error) {
	salt := os.Getenv("SALT")
	pass := salt + password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		logger.Print(logger.ERROR, "パスワードのハッシュ化中にエラーが発生しました: %v", err)
		return "", err
	}
	encodedPassword := string(hashedPassword)
	return encodedPassword, nil
}

// ヘッダーのトークンと認証情報をチェック
func CheckHeader(res http.ResponseWriter, req *http.Request) *ctm_err.CustomError {
	auth := req.Header.Get("authorization")
	tkn := req.Header.Get("token")

	// tknのチェックをスキップするエンドポイントを定義
	skipEndpoints := []string{
		"/api/v1/permission",
		"/api/v1/login",
		"/api/v1/token/create",
	}

	if !contains(skipEndpoints, req.URL.Path) {
		logger.Print(logger.INFO, "トークンチェック対象外のためスキップします. エンドポイント: %s", req.URL.Path)
		return nil
	}

	tokenInfo, ok := cache.GetCache[session.TokenInfo](tkn, false)
	if !ok {
		err := ctm_err.NewRequestErr("トークンが不正です")
		logger.Print(logger.ERROR, "トークンが不正です")
		return err
	}
	//クライアントIDをもとにauthorizationが想定通りの設定内容かチェック
	if auth != fmt.Sprintf("Bearer %s", tokenInfo.ClientId) {
		err := ctm_err.NewRequestErr("認証情報が不正です")
		logger.Print(logger.ERROR, "認証情報が不正です")
		return err
	}

	return nil
}

// contains checks if a string is present in a slice of strings
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
