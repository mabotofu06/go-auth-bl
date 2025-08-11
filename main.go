package main

import (
	"fmt"
	"go-auth-bl/cache"
	con "go-auth-bl/internal/controller"
	"go-auth-bl/internal/def"
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	if err := cache.Init(); err != nil {
		fmt.Printf("キャッシュ初期化エラー: %v\n", err)
		return
	}

	//認可コード要求API
	http.HandleFunc(def.CONFIG.API["permission"].Endpoint, ApiWrapper(con.GetPermission))
	// ログインAPI
	//curl -X POST http://localhost/api/login -H "Content-Type: application/json" -d "{\"userId\": \"elf_hinmel\", \"email\": \"\", \"password\": \"password\"}"
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
	http.HandleFunc("/", ApiWrapper(root))

	fmt.Println("Starting server at port 8080")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("sesid", "session-id-1234-5678-9000")
	w.Header().Set("acc", "acc-id-1234-5678-9000")

	http.FileServer(http.Dir("./build")).ServeHTTP(w, r)
	return
}

func ApiWrapper(
	controller func(w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	return controller
}
