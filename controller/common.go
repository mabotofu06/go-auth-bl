package controller

import (
	"encoding/json"
	"fmt"
	apiif "go-auth-bl/dto/if"
	a_err "go-auth-bl/error"
	"go-auth-bl/middleware"
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
func ReqMethodCheck(res http.ResponseWriter, req *http.Request, method string) *a_err.CustomError {
	if req.Method != method {
		err := a_err.NewRequestErr("リクエストメソッドが不適切です")

		fmt.Printf("リクエストメソッドが不適切です: %s\n", req.Method)
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
		middleware.ResError(res, err)
		return nil, err
	}

	defer req.Body.Close()
	var request T

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		err := a_err.NewRequestErr("リクエストボディが不適切です")
		fmt.Println("リクエストボディエンコード中にエラーが発生しました:", err)
		middleware.ResError(res, err)
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
