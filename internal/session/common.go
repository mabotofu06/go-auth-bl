package session

import "net/http"

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

func SetSessionId(res http.ResponseWriter, sessionId string) {
	http.SetCookie(res, &http.Cookie{
		Name:     "sesid",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   30 * 60,
		HttpOnly: true,
		Secure:   false, // HTTPSならtrue
		SameSite: http.SameSiteLaxMode,
	})
}
