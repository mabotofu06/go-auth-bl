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
	ReqMethodCheck(req, GET)

	queryParams := req.URL.Query()
	code := queryParams.Get("code")
	//redirect_uri := queryParams.Get("redirect_uri")
	fmt.Printf("認可コード取得処理を開始します code=%s\n", code)

	// セッション情報発行
	// 情報をキャッシュに保存
	// セッションIDをクッキーに保存
	// UIを返却（ログイン画面）
}
