package main

import (
	"fmt"
	con "go-auth-bl/controller"
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	//認可コード要求API
	http.HandleFunc("/api/v1/auth_permission", ApiWrapper(con.GetPermission))

	// ログインAPI
	//curl -X POST http://localhost/api/login -H "Content-Type: application/json" -d "{\"userId\": \"elf_hinmel\", \"email\": \"\", \"password\": \"password\"}"
	http.HandleFunc("/api/v1/login", ApiWrapper(con.PostLogin))

	//http://localhost/ にアクセスすると画面が返却
	http.HandleFunc("/", ApiWrapper(root))

	fmt.Println("Starting server at port 8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("sesid", "session-id-1234-5678-9000")
	w.Header().Set("acc", "acc-id-1234-5678-9000")
	w.Header().Set("sesid", "session-id-1234-5678-9000")
	w.Header().Set("sesid", "session-id-1234-5678-9000")

	http.FileServer(http.Dir("./build")).ServeHTTP(w, r)
	return
	// if r.URL.Path != "/" {
	// 		// 静的ファイルを返す
	// 		http.FileServer(http.Dir("./build")).ServeHTTP(w, r)
	// 		return
	// 	}

	// 	con.GetPermission(w, r)
}

func ApiWrapper(
	controller func(w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	return controller
}
