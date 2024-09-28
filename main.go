package main

import (
	"encoding/json"
	"fmt"
	"go-auth-bl/controller"
	"go-auth-bl/service"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received!!")
	fmt.Println("Jsonデータを返します")

	res := map[string]string{
		"message": "Hello, World!",
	}

	userId := "elf_hinmel"

	service.GetUserAuthByUserId(userId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: nil, // DefaultServeMux を使用
	}

	//http://localhost/ にアクセスするとHello, World!が表示される
	http.HandleFunc("/", root)
	// ログインAPI
	//curl -X POST http://localhost/api/login -H "Content-Type: application/json" -d "{\"userId\": \"elf_hinmel\", \"email\": \"\", \"password\": \"password\"}"
	http.HandleFunc("/api/login", controller.PostLogin)

	fmt.Println("Starting server at port 8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
