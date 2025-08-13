package repository

// go:generate mockgen -source=client_info.go -destination=mock/mock_client_info_repository.go -package=repository_mock
import (
	"database/sql"
	"fmt"
	"go-auth-bl/internal/dto"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
)

type IClientInfoRepository interface {
	GetClientInfoByClientId(clientId string, db *sql.DB) (*dto.ClientInfo, *ctm_err.CustomError)
}
type clientInfoRepository struct{}

var ClientInfoRepository IClientInfoRepository = &clientInfoRepository{}

// クライアントIDを元にクライアント情報を取得(プライマリーキーを元に検索のため1件のみ取得)
func (clientInfoRepository) GetClientInfoByClientId(clientId string, db *sql.DB) (*dto.ClientInfo, *ctm_err.CustomError) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_flag = 0 AND client_id = $1", MST_CLIENT_INFO)
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
