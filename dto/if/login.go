package apiif

type ReqLogin struct {
	UsrId    string `json:"userId"`
	Password string `json:"password"`
}

type ResLogin struct {
	UsrId   string `json:"userId"`
	Session string `json:"session"`
}
