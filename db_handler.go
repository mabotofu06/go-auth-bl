package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-auth-bl/dto"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() {
	// データベースに接続
	db, err := sql.Open(
		"postgres",
		"host=go-auth-db port=5432 user=go-auth-db password=postgres dbname=go-auth-db sslmode=disable",
	)
	ErrorCheck(err)
	defer db.Close() // 関数終了時に接続を閉じる

	// 接続を確認
	err = db.Ping()
	ErrorCheck(err)

	fmt.Println("Successfully connected to the database!")

	// クエリの実行例
	rows, err := db.Query("SELECT * FROM mng_user_auth_tbl")
	ErrorCheck(err)
	defer rows.Close() // 関数終了時にクエリを閉じる

	//var userAuthList []dto.UserAuth
	for rows.Next() {
		var userAuth dto.UserAuth

		// クエリの結果をDTOにマッピング
		err := rows.Scan(
			&userAuth.UserId,
			&userAuth.Password,
			&userAuth.SessionToken,
			&userAuth.LastSessinonTime,
			&userAuth.DeleteFlag,
			&userAuth.CreateDateTime,
			&userAuth.UpdateDateTime,
			&userAuth.DeleteDate,
		)
		ErrorCheck(err)

		// DTOをJSON形式に変換
		userJSON, err := json.Marshal(userAuth)
		if err != nil {
			log.Fatal(err)
		}

		// JSON形式のDTOを出力
		fmt.Println(string(userJSON))

		//userAuthList = append(userAuthList, userAuth)
	}

	ErrorCheck(err)

	//	return &userList, nil
}

func ErrorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
