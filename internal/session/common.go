package session

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

type PermissionInfo struct {
	ClientId    string
	RedirectUri string
	Scope       string
	State       string
}

type TokenInfo struct {
	AccessToken string
	RedirectUri string
}

var Session = sessions.NewCookieStore([]byte("go-auth-session"))

func init() {
	// gobパッケージに構造体を登録
	gob.Register(PermissionInfo{})
	gob.Register(TokenInfo{})
}

func Get(r *http.Request) (*sessions.Session, error) {
	return Session.Get(r, "auth")
}

func GetValue(key string, r *http.Request) interface{} {
	session, _ := Get(r)
	return session.Values[key]
}

func SetValue(key string, value interface{}, expireSec int, session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	session.Values[key] = value
	session.Options = &sessions.Options{
		MaxAge: expireSec,
	}
	session.Save(r, w)
}
