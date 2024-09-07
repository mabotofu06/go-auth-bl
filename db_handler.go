package main

import (
	"database/sql"
	"fmt"
	"log"
)

func connect_db() {
	// データベース接続情報
	connStr := "user=username password=password dbname=mydb sslmode=disable"

	// データベースに接続
	db, err := sql.Open("postgres", connStr)
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

	id := 1

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
