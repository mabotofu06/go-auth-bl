package controller

import (
	"encoding/json"
	"go-auth-bl/service"
	"net/http"
)

type LoginRequest struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func PostLogin(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(
			res,
			"Invalid request method",
			http.StatusMethodNotAllowed,
		)
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
		http.Error(
			res,
			"Error fetching user auth data",
			http.StatusInternalServerError,
		)

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
