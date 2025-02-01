package apiif

type Header struct {
	ReqId  string `json:"requestId"`
	SessId string `json:"sessionId"`
	UsrId  string `json:"userId"`
}

type Response[T any] struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Msg    string `json:"message"`
	Data   *T     `json:"data"`
}
