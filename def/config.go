package def

type GoAuthApi struct {
	name     string
	method   string
	endpoint string
}

type GoAuthConfig struct {
	PORT int
	API  map[string]GoAuthApi
}

var CONFIG = GoAuthConfig{
	PORT: 3000,
	API: map[string]GoAuthApi{
		"auth_permission": {name: "認可ID発行", method: "GET", endpoint: "/api/auth_permission"},
		"login":           {name: "ログイン", method: "POST", endpoint: "/api/login"},
		"access_token":    {name: "アクセストークン要求", method: "GET", endpoint: "/api/access_token"},
	},
}

var ERROR_MESSAGE = map[string]string{
	"W0001": "リクエストメソッドが不適切です",
	"W0002": "リクエストパラメータが不適切です",
	"W0003": "認証エラーが発生しました",
	"W0005": "ユーザー名またはパスワードが違います",

	"E0001": "予期せぬエラーが発生しました",
}
