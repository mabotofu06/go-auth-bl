package service

import (
	"database/sql"
	"fmt"
	"go-auth-bl/dto"
	"go-auth-bl/repository"
	"log"

	_ "github.com/lib/pq"
)

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {
	// データベースに接続
	dbConnectInfo := "host=go-auth-db port=5432 user=go-auth-db password=postgres dbname=go-auth-db sslmode=disable"

	db, err := sql.Open("postgres", dbConnectInfo)
	errorCheck(err)

	// 接続を確認
	err = db.Ping()
	errorCheck(err)

	fmt.Println("Successfully connected to the database!")
	return db
}

// Coution: 大文字でないと外部パッケージから参照できない
// ユーザIDを元にユーザ認証情報を取得
func GetUserAuthByUserId(userId string) (*dto.UserAuth, error) {
	db := ConnectDB()
	defer db.Close() // 関数終了時に接続を閉じる

	// サービス層を呼び出してデータを取得
	userAuths, err := repository.GetUserAuthsByUserId(userId, db)
	errorCheck(err)

	fmt.Printf("userAuths: %v\n", userAuths)

	if len(userAuths) == 0 {
		fmt.Println("No user auth data found")
		return nil, nil
	}

	return &userAuths[0], nil
}
