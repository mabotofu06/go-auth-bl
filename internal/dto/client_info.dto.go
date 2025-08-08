package dto

//クライアント情報を表す構造体
type ClientInfo struct {
	ClientId   string `json:"clientId"`
	ClientName string `json:"clientName"`
	ClientHost string `json:"clientHost"`

	DtoBase
}
