package repository

import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
	"log"
)

// Coution: 大文字でないと外部パッケージから参照できない
func GetClientInfoByClientId(clientId string, db *sql.DB) (*dto.ClientInfo, error) {
	const tableName = "mng_client_mst"
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND client_id = $1", tableName)
	// クライアントIDを元にクライアント情報を取得(プライマリーキーを元に検索のため1件のみ取得)
	row := db.QueryRow(query, clientId)
	if row == nil {
		return nil, fmt.Errorf("クライアント情報が見つかりません: %s", clientId)
	}

	var clientInfo dto.ClientInfo
	if err := row.Scan(
		&clientInfo.ClientId,
		&clientInfo.ClientName,
		&clientInfo.ClientHost,
		&clientInfo.DeleteFlag,
		&clientInfo.CreateDateTime,
		&clientInfo.UpdateDateTime,
		&clientInfo.DeleteDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Fatal(err)
	}

	return &clientInfo, nil
}
