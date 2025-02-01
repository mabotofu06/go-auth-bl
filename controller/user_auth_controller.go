package controller

import (
	"encoding/json"
	"fmt"
	apiif "go-auth-bl/dto/if"
	a_err "go-auth-bl/error"
	"go-auth-bl/service"
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

func IsMethod(method string, req *http.Request) bool {
	return req.Method == method
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

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func PostLogin(res http.ResponseWriter, req *http.Request) {
	// POSTメソッド以外はエラー
	if !IsMethod(POST, req) {
		panic(a_err.BadRequestErr)
	}
	// リクエストボディをJSONとしてパース
	var request apiif.ReqLogin
	var err error
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		panic(a_err.BadRequestErr)
	}
	fmt.Printf("Login request: %+v\n", request)
	// サービス層を呼び出してデータを取得
	userAuth, err := service.GetUserAuthByUserId(request.UsrId)
	if err != nil {
		if err == a_err.NotFoundErr {
			panic(a_err.NotFoundErr)
		}
		panic(a_err.InternalServerErr)
	}

	// パスワードが一致するか確認
	EncodePassword(request.Password)
	passCheck, err := service.PasswordCheck(userAuth, os.Getenv("SALT")+request.Password)
	if err != nil {
		panic(a_err.InternalServerErr)
	}
	if !passCheck {
		panic(a_err.NotFoundErr)
	}

	fmt.Printf("ユーザ情報: %v\n", userAuth)
	// レスポンスを返す
	res.Header().Set("Content-Type", "application/json")
	// res.Header().Set("Access-Control-Allow-Origin", "*")
	// res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// res.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	body := apiif.Response[apiif.ResLogin]{
		Status: http.StatusOK,
		Code:   "I0001",
		Msg:    "Success",
		Data: &apiif.ResLogin{
			UsrId:   userAuth.UserId,
			Session: "dummy_session",
		},
	}
	json, _ := json.Marshal(body)
	res.WriteHeader(http.StatusOK)
	res.Write(json)
}

func ResError(res http.ResponseWriter, status int, code string, msg string) {
	res.Header().Set("Content-Type", "application/json")

	body := apiif.Response[interface{}]{
		Status: status,
		Code:   code,
		Msg:    msg,
		Data:   nil,
	}
	json, _ := json.Marshal(body)

	res.WriteHeader(status)
	res.Write(json)
}
