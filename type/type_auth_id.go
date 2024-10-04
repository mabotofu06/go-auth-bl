package go_auth_type

type AuthIdApiRequest struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
