package dto

type UserInfo struct {
	UserId   string  `json:"userId"`
	UserName string  `json:"userName"`
	Email    *string `json:"email"`
	Gender   *string `json:"gender"`
	Age      *string `json:"age"`

	DtoBase
}
