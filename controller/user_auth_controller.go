package controller

import (
	"encoding/json"
	"fmt"
	"go-auth-bl/service"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func IsMethod(method string, req *http.Request) bool {
	return req.Method == method
}

func SendNotAllowMethod(req *http.Request, res http.ResponseWriter) {
	fmt.Println("Invalid request method for", req.URL.Path)
	http.Error(
		res,
		"Invalid request method",
		http.StatusMethodNotAllowed,
	)
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
	if !IsMethod(POST, req) {
		SendNotAllowMethod(req, res)
		return
	}

	// リクエストボディをJSONとしてパース
	var loginRequest LoginRequest
	var err error

	err = json.NewDecoder(req.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(
			res,
			"Error parsing JSON",
			http.StatusBadRequest,
		)

		return
	}

	// サービス層を呼び出してデータを取得
	userAuth, err := service.GetUserAuthByUserId(loginRequest.UserId)
	if err != nil {
		fmt.Println("Error fetching user auth data:", err)
		http.Error(res, "Error fetching user auth data", http.StatusInternalServerError)
		return
	}

	EncodePassword(loginRequest.Password)

	// パスワードが一致するか確認
	if !service.PasswordCheck(userAuth, os.Getenv("SALT")+loginRequest.Password) {
		fmt.Println("Password does not match")
		http.Error(res, "Password does not match", http.StatusUnauthorized)
		return
	}

	// レスポンスとしてJSONを返す
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	err = json.NewEncoder(res).Encode(userAuth)
	if err != nil {
		http.Error(
			res,
			"Error encoding JSON",
			http.StatusInternalServerError,
		)
	}
}
