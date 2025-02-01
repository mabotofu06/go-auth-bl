package common_db

import (
	"database/sql"
	"fmt"
	a_err "go-auth-bl/error"
	"os"

	_ "github.com/lib/pq"
)

// データベースに接続
func ConnectDB() *sql.DB {
	dbConnectInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sql.Open("postgres", dbConnectInfo)
	if err != nil {
		a_err.Throw(a_err.NewDBErr("予期せぬエラーが発生しました"))
	}
	// 接続を確認
	err = db.Ping()
	if err != nil {
		a_err.Throw(a_err.NewDBErr("予期せぬエラーが発生しました"))
	}
	fmt.Println("Successfully connected to the database!")
	return db
}
