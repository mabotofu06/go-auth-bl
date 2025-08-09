package controller

import (
	"encoding/json"
	"fmt"
	"go-auth-bl/cache"
	apiif "go-auth-bl/internal/dto/if"
	"go-auth-bl/internal/middleware"
	"go-auth-bl/internal/session"
	a_err "go-auth-bl/pkg/error"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var Store = sessions.NewCookieStore([]byte("go-auth-session"))

// リクエストメソッドが不適切な場合はエラーを返す
func ReqMethodCheck(res http.ResponseWriter, req *http.Request, method string) *a_err.CustomError {
	if req.Method != method {
		err := a_err.NewRequestErr("リクエストメソッドが不適切です")

		fmt.Printf("リクエストメソッドが不適切です: %s\n", req.Method)
		//FIXME:各処理でのResError処理は各controller内部で呼ばれるように修正
		middleware.ResError(res, err)
		return err
	}
	return nil
}

// POSTリクエストのリクエストボディを取得（リクエストボディが不適切な場合はエラーを返す）
func GetReqBody[T any](res http.ResponseWriter, req *http.Request) (*T, *a_err.CustomError) {
	if req.Body == nil {
		err := a_err.NewRequestErr("リクエストボディが空です")
		fmt.Println("リクエストボディが空です")
		//FIXME:各処理でのResError処理は各controller内部で呼ばれるように修正
		middleware.ResError(res, err)
		return nil, err
	}

	defer req.Body.Close()
	var request T

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		err := a_err.NewRequestErr("リクエストボディが不適切です")
		fmt.Println("リクエストボディエンコード中にエラーが発生しました:", err)
		//FIXME:各処理でのResError処理は各controller内部で呼ばれるように修正
		middleware.ResError(res, err)
		return nil, err
	}

	fmt.Printf("request: %+v\n", request)
	return &request, nil
}

// APIの正常終了時のレスポンスを返す
func ResOk[T any](res http.ResponseWriter, data *T) {
	fmt.Printf("response data: %+v\n", data)

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
		//FIXME:各処理でのResError処理は各controller内部で呼ばれるように修正
		middleware.ResError(res, a_err.NewServerErr("予期せぬエラーが発生しました"))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write(json)
}

// bcryptを使ってパスワードをハッシュ化
// エンコードされた文字列の長さは60文字
func EncodePassword(password string) (string, error) {
	salt := os.Getenv("SALT")
	pass := salt + password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return "", err
	}
	encodedPassword := string(hashedPassword)
	// fmt.Println("Encoded password:", encodedPassword)
	return encodedPassword, nil
}

// ヘッダーのトークンと認証情報をチェック
func CheckHeader(res http.ResponseWriter, req *http.Request) *a_err.CustomError {
	auth := req.Header.Get("authorization")
	tkn := req.Header.Get("token")

	// tknのチェックをスキップするエンドポイントを定義
	skipEndpoints := []string{
		"/api/v1/permission",
		"/api/v1/login",
		"/api/v1/token/create",
	}

	if !contains(skipEndpoints, req.URL.Path) {
		fmt.Printf("トークンチェック対象外のためスキップします. エンドポイント: %s\n", req.URL.Path)
		return nil
	}

	tokenInfo, ok := cache.GetCache[session.TokenInfo](tkn, false)
	if !ok {
		err := a_err.NewRequestErr("トークンが不正です")
		fmt.Println("トークンが不正です")
		return err
	}
	//クライアントIDをもとにauthorizationが想定通りの設定内容かチェック
	if auth != fmt.Sprintf("Bearer %s", tokenInfo.ClientId) {
		err := a_err.NewRequestErr("認証情報が不正です")
		fmt.Println("認証情報が不正です")
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
