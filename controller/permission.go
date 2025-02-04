package controller

import (
	"fmt"
	"net/http"
)

/**
* @param w http.ResponseWriter
* @param r *http.Request
 */
func GetPermission(res http.ResponseWriter, req *http.Request) {
	if ReqMethodCheck(res, req, GET) != nil {
		return
	}

	queryParams := req.URL.Query()
	code := queryParams.Get("code")
	//redirect_uri := queryParams.Get("redirect_uri")
	fmt.Printf("認可コード発行処理を開始します code=%s\n", code)

	// Basic認証
	req.Header.Get("Authorization")

	// セッション情報発行
	// 情報をキャッシュに保存
	// セッションIDをクッキーに保存
	// UIを返却（ログイン画面）
	http.FileServer(http.Dir("./build")).ServeHTTP(res, req)
}
