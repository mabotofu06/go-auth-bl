package dto

type UserInfo struct {
	UserId   string  `json:"userId"`
	UserName string  `json:"userName"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Gender   *string `json:"gender"`

	DtoBase
}
