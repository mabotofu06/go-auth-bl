package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	id   int
	name string
}

// const dbHost = "go-auth-db"
// const port = "5432"
// const dbUser = "go-auth-db"
// const dbPassword = "postgres"
// const dbName = "go-auth-db"

func ConnectDB() (*[]User, error) {
	// データベースに接続
	db, err := sql.Open(
		"postgres",
		"host=go-auth-db port=5432 user=go-auth-db password=postgres dbname=go-auth-db sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close() // 関数終了時に接続を閉じる

	// 接続を確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")

	userList := []User{} // ユーザ情報を格納するスライス

	id := 3

	// クエリの実行例
	rows, err := db.Query(
		"SELECT id, name FROM users WHERE id = $1",
		id,
	)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close() // 関数終了時にクエリを閉じる

	for rows.Next() {
		var user User
		var id int
		var name string
		if err := rows.Scan(&user.id, &user.name); err != nil {
			log.Fatal(err)
			return nil, err
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
		userList = append(userList, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &userList, nil
}
