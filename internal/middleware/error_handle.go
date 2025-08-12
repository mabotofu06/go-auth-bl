package middleware

import (
	"encoding/json"
	apiif "go-auth-bl/internal/dto/if"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
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
				customErr, ok := err.(*ctm_err.CustomError)
				if !ok {
					customErr = ctm_err.NewServerErr("予期せぬエラーが発生しました")
				}
				ResError(w, customErr)
			}()

			next.ServeHTTP(w, r)
		})
}

func ResError(res http.ResponseWriter, err *ctm_err.CustomError) {
	res.Header().Set("Content-Type", "application/json")

	if err == nil {
		logger.Print(logger.ERROR, "エラーがnilです")
		err = ctm_err.NewServerErr("予期せぬエラーが発生しました")
	}

	body := apiif.Response[any]{
		Status: err.Status,
		Code:   err.Code,
		Type:   err.Type,
		Msg:    err.Msg,
		Data:   nil,
	}
	logger.Print(logger.DEBUG, "response : %+v", body)
	json, _ := json.Marshal(body)

	res.WriteHeader(err.Status)
	res.Write(json)
}
