package main

import (
	"fmt"
	con "go-auth-bl/controller"
	"net/http"
)

// func root() {
// 	fmt.Println("ログインページを返します。")

// 	// HTMLファイルを返す
// 	return http.FileServer(http.Dir("./build"))
// }

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: nil, // DefaultServeMux を使用
	}

	// ログインAPI
	//curl -X POST http://localhost/api/login -H "Content-Type: application/json" -d "{\"userId\": \"elf_hinmel\", \"email\": \"\", \"password\": \"password\"}"
	http.HandleFunc("/api/login", con.PostLogin)

	//http://localhost/ にアクセスするとHello, World!が表示される
	http.Handle("/", http.FileServer(http.Dir("./build")))

	fmt.Println("Starting server at port 8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
