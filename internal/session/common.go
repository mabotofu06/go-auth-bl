package session

type PermissionInfo struct {
	ClientId    string
	RedirectUri string
	Scope       string
	State       string
}

type CodeInfo struct {
	AccessToken string
	UserId      string
	ClientId    string
	Scope       string
	RedirectUri string
}

type TokenInfo struct {
	UserId   string
	Scope    string
	ClientId string
}
