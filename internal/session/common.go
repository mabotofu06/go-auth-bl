package session

import (
	"encoding/gob"
	"fmt"
	a_err "go-auth-bl/pkg/error"
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

func getAuthSession(r *http.Request) (*sessions.Session, error) {
	return Session.Get(r, "auth")
}

func GetValue[T any](key string, r *http.Request) (*T, *a_err.CustomError) {
	session, err := getAuthSession(r)
	if err != nil {
		fmt.Printf("認証セッション取得時にエラーが発生しました: %v\n", err)
		return nil, a_err.NewAuthErr("認証エラー")
	}

	val, ok := session.Values[key]
	if !ok {
		fmt.Printf("認証セッション取得時にエラーが発生しました: %v\n", err)
		return nil, a_err.NewAuthErr("認証エラー")
	}
	value, ok := val.(T)
	if !ok {
		fmt.Printf("認証セッションの型が一致しません: %v\n", val)
		return nil, a_err.NewAuthErr("認証エラー")
	}

	fmt.Printf("認証セッションから値を取得しました key: %s value: %v\n", key, value)

	return &value, nil
}

func SetValue[T any](w http.ResponseWriter, r *http.Request, key string, value T, expireSec int) *a_err.CustomError {
	session, err := getAuthSession(r)
	if err != nil {
		fmt.Printf("認証セッション設定にエラーが発生しました: %v\n", err)
		return a_err.NewAuthErr("認証エラー")
	}

	switch any(value).(type) {
	case PermissionInfo, TokenInfo:
		break
	default:
		fmt.Printf("認証セッションの型が異なります")
		return a_err.NewServerErr("認証エラー")
	}

	session.Values[key] = value
	session.Options = &sessions.Options{
		MaxAge: expireSec,
	}
	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("認証セッションの保存にエラーが発生しました: %v\n", err)
		return a_err.NewServerErr("セッション保存エラー")
	}
	fmt.Printf("認証セッションに値を設定しました key: %s value: %v\n", key, value)
	return nil
}
