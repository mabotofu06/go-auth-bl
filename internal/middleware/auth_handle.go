package middleware

import (
	"fmt"
	a_err "go-auth-bl/pkg/error"
	"net/http"
	"strings"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				fmt.Println("APIのエンドポイントでないため認証をスキップします")
				next.ServeHTTP(w, r)
				return
			}

			if r.URL.Path == "/api/login" {
				fmt.Println("ログインAPIのエンドポイントであるため認証をスキップします")
				next.ServeHTTP(w, r)
				return
			}

			//sesid := r.Header.Get("sessionId")
			//TODO:セッションチェック
			a_err.Throw(a_err.NewAuthErr("認証中にエラーが発生しました"))
		})
}
