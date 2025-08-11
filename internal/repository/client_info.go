package repository

import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
)

// Coution: 大文字でないと外部パッケージから参照できない
func GetClientInfoByClientId(clientId string, db *sql.DB) (*dto.ClientInfo, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND client_id = $1", MST_CLIENT_INFO)
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
		fmt.Printf("db.Query: %v\n", err)
		return nil, err
	}

	return &clientInfo, nil
}
