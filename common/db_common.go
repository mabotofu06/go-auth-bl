package common_db

import (
	"database/sql"
	"fmt"
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
	ErrorCheck(err)

	// 接続を確認
	err = db.Ping()
	ErrorCheck(err)

	fmt.Println("Successfully connected to the database!")
	return db
}
