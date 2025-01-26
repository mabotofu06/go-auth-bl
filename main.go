package main

import (
	"fmt"
	con "go-auth-bl/controller"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	//TODO:認証チェック

	w.Header().Set("sesid", "session-id-1234-5678-9000")
	w.Header().Set("acc", "acc-id-1234-5678-9000")
	w.Header().Set("sesid", "session-id-1234-5678-9000")
	w.Header().Set("sesid", "session-id-1234-5678-9000")

	// 静的ファイルを返す
	http.FileServer(http.Dir("./build")).ServeHTTP(w, r)

	if r.URL.Path != "/" {
		return
	}

	queryParams := r.URL.Query()
	for key, values := range queryParams {
		for _, value := range values {
			fmt.Printf("Query parameter: %s = %s\n", key, value)
		}
	}

	fmt.Printf("認可コード取得処理を開始します code=%s\n", queryParams["code"])
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: nil, // DefaultServeMux を使用
	}

	// ログインAPI
	//curl -X POST http://localhost/api/login -H "Content-Type: application/json" -d "{\"userId\": \"elf_hinmel\", \"email\": \"\", \"password\": \"password\"}"
	http.HandleFunc("/api/login", con.PostLogin)

	//http://localhost/ にアクセスするとHello, World!が表示される
	http.HandleFunc("/", root)

	fmt.Println("Starting server at port 8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
