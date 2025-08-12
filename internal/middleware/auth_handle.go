package middleware

import (
	a_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/http"
	"strings"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				logger.Print(logger.DEBUG, "APIのエンドポイントでないため認証をスキップします")
				next.ServeHTTP(w, r)
				return
			}

			if r.URL.Path == "/api/login" {
				logger.Print(logger.DEBUG, "ログインAPIのエンドポイントであるため認証をスキップします")
				next.ServeHTTP(w, r)
				return
			}

			//sesid := r.Header.Get("sessionId")
			//TODO:セッションチェック
			a_err.Throw(a_err.NewAuthErr("認証中にエラーが発生しました"))
		})
}
