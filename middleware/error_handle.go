package middleware

import (
	"go-auth-bl/controller"
	a_err "go-auth-bl/error"
	"log"
	"net/http"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("cacthed err: %v", err)
				if customErr, ok := err.(*a_err.CustomError); ok {
					controller.ResError(w, customErr.Status, customErr.Code, customErr.Msg)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
