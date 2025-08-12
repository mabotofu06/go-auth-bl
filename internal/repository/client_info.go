package repository

import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
)

// Coution: 大文字でないと外部パッケージから参照できない
func GetClientInfoByClientId(clientId string, db *sql.DB) (*dto.ClientInfo, *ctm_err.CustomError) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND client_id = $1", MST_CLIENT_INFO)
	// クライアントIDを元にクライアント情報を取得(プライマリーキーを元に検索のため1件のみ取得)
	row := db.QueryRow(query, clientId)
	if row == nil {
		return nil, ctm_err.NewNotFoundErr("クライアント情報が見つかりませんでした")
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
		logger.Print(logger.ERROR, "db.Query: %v", err)
		return nil, ctm_err.UnexpectedDBErr
	}

	return &clientInfo, nil
}
