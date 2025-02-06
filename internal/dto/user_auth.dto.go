package dto

type UserAuth struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
	//null許容の場合はポインタ型にする
	PasswordHistory1 *string `json:"passwordHistory1"`
	PasswordHistory2 *string `json:"passwordHistory2"`
	PasswordHistory3 *string `json:"passwordHistory3"`
	PasswordFailCnt  int     `json:"passwordFailCnt"`
	PasswordLockFlag int     `json:"passwordLockFlag"`

	DtoBase
}

type DtoBase struct {
	DeleteFlag     string  `json:"deleteFlag"`
	CreateDateTime string  `json:"createDateTime"`
	UpdateDateTime *string `json:"updateDateTime"`
	DeleteDate     *string `json:"deleteDate"`
}
