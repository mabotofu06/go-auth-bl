package dto

type UserAuth struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
	//null許容の場合はポインタ型にする
	SessionToken     *string `json:"sessionToken"`
	LastSessinonTime *string `json:"lastSessinonTime"`

	DtoBase
}

type DtoBase struct {
	DeleteFlag     string  `json:"deleteFlag"`
	CreateDateTime string  `json:"createDateTime"`
	UpdateDateTime *string `json:"updateDateTime"`
	DeleteDate     *string `json:"deleteDate"`
}
