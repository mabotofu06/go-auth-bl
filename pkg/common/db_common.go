package common_db

import (
	"database/sql"
	"fmt"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"os"

	_ "github.com/lib/pq"
)

// データベースに接続
func ConnectDB() (*sql.DB, *ctm_err.CustomError) {
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
		return nil, ctm_err.UnexpectedDBErr
	}
	// 接続を確認
	err = db.Ping()
	if err != nil {
		return nil, ctm_err.UnexpectedDBErr
	}
	logger.Print(logger.DEBUG, "DB接続情報: %s", dbConnectInfo)
	logger.Print(logger.INFO, "データベースとの接続に成功しました")
	return db, nil
}
