package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONデータを表す構造体
type Data struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	//http://localhost/ にアクセスするとHello, World!が表示される
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received!!")
		fmt.Println("Jsonデータを返します")

		res := map[string]string{
			"message": "Hello, World!",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		}
	})

	//http://localhost/data にPOSTリクエストを送ると、リクエストボディをそのまま返す
	//curl -X POST http://localhost/data -H "Content-Type: application/json" -d "{\"name\": \"John Doe\", \"email\": \"john.doe@example.com\"}"
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// リクエストボディをJSONとしてパース
		var data Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		// レスポンスとしてJSONを返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		}
	})

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// リクエストボディをJSONとしてパース
		var data Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		//DBにアクセス
		connect_db()
	})

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
