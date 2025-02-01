package main

import (
	"fmt"
	con "go-auth-bl/controller"
	"go-auth-bl/middleware"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	//TODO:認証チェック
	//s := r.Header.Get("sesid")
	// if s != "auth" {
	// 	fmt.Println("認証エラー")
	// 	http.Error(w, "認証エラー", http.StatusUnauthorized)
	// 	return
	// }

	w.Header().Set("sesid", "session-id-1234-5678-9000")
	w.Header().Set("acc", "acc-id-1234-5678-9000")
	w.Header().Set("sesid", "session-id-1234-5678-9000")
	w.Header().Set("sesid", "session-id-1234-5678-9000")

	// 静的ファイルを返す
	http.FileServer(http.Dir("./build")).ServeHTTP(w, r)

	if r.URL.Path != "/" {
		return
	}

	//fmt.Printf("認可コード取得処理を開始します code=%s\n", queryParams["code"])
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.ErrorHandler(http.DefaultServeMux),
	}

	//http://localhost/ にアクセスすると画面が返却
	http.HandleFunc("/", root)
	// ログインAPI
	//curl -X POST http://localhost/api/login -H "Content-Type: application/json" -d "{\"userId\": \"elf_hinmel\", \"email\": \"\", \"password\": \"password\"}"
	http.HandleFunc("/api/login", con.PostLogin)

	fmt.Println("Starting server at port 8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
