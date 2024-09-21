package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() {
	// データベースに接続
	db, err := sql.Open(
		"postgres",
		"host=127.0.0.1 port=5432 user=go-auth-db password=postgres dbname=go-auth-db sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 接続を確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")

	id := 3

	// クエリの実行例
	rows, err := db.Query(
		"SELECT id, name FROM users WHERE id = $1",
		id,
	)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
