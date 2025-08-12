package main

import (
	"go-auth-bl/internal/cache"
	con "go-auth-bl/internal/controller"
	"go-auth-bl/internal/def"
	"go-auth-bl/pkg/logger"
	"net/http"
	"os"
)

func main() {
	port := ":8080"
	server := http.Server{
		Addr:    port,
		Handler: http.DefaultServeMux,
	}
	//ロガー初期化(環境変数を元に出力するログレベルを設定)
	logger.Init(os.Getenv("LOG_LEVEL"))
	//キャッシュ初期化
	if err := cache.Init(); err != nil {
		logger.Print(logger.ERROR, "キャッシュ初期化エラー: %v", err)
		return
	}

	//各APIのエンドポイントを設定
	//認可コード要求API
	http.HandleFunc(def.CONFIG.API["permission"].Endpoint, ApiWrapper(con.GetPermission))
	// ログインAPI
	http.HandleFunc(def.CONFIG.API["login"].Endpoint, ApiWrapper(con.PostLogin))
	//アクセストークン要求API
	http.HandleFunc(def.CONFIG.API["get_token"].Endpoint, ApiWrapper(con.GetAccessToken))
	//トークン検証API
	http.HandleFunc(def.CONFIG.API["check_token"].Endpoint, ApiWrapper(con.GetValidToken))
	//トークン削除API
	http.HandleFunc(def.CONFIG.API["delete_token"].Endpoint, ApiWrapper(con.DeleteToken))
	//ユーザ登録API
	http.HandleFunc(def.CONFIG.API["create_user"].Endpoint, ApiWrapper(con.PostCreateUser))
	//ユーザ情報取得API
	http.HandleFunc(def.CONFIG.API["get_user"].Endpoint, ApiWrapper(con.GetUserInfo))

	//http://localhost/ にアクセスすると画面が返却
	http.HandleFunc("/", root)

	logger.Print(logger.INFO, "サーバー起動中: ポート %s", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Print(logger.ERROR, "サーバー起動エラー: %v", err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./build")).ServeHTTP(w, r)
	return
}

func ApiWrapper(
	controller func(w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	//logger.Print(logger.INFO, "API: %s 呼び出します", apiInfo.Name)
	return controller
}
