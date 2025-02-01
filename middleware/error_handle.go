package middleware

import (
	"encoding/json"
	apiif "go-auth-bl/dto/if"
	a_err "go-auth-bl/error"
	"log"
	"net/http"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err == nil {
					next.ServeHTTP(w, r)
					return
				}

				log.Printf("cacthed err: %v", err)
				customErr, ok := err.(*a_err.CustomError)
				if !ok {
					customErr = a_err.NewServerErr("予期せぬエラーが発生しました")
				}
				resError(w, customErr)
			}()

			next.ServeHTTP(w, r)
		})
}

func resError(res http.ResponseWriter, err *a_err.CustomError) {
	res.Header().Set("Content-Type", "application/json")

	body := apiif.Response[any]{
		Status: err.Status,
		Code:   err.Code,
		Type:   err.Type,
		Msg:    err.Msg,
		Data:   nil,
	}
	json, _ := json.Marshal(body)

	res.WriteHeader(err.Status)
	res.Write(json)
}
