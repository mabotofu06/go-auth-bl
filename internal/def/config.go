package def

type GoAuthApi struct {
	Id       string
	Name     string
	Method   string
	Endpoint string
}

type GoAuthConfig struct {
	PORT int
	API  map[string]GoAuthApi
}

var CONFIG = GoAuthConfig{
	PORT: 3000,
	API: map[string]GoAuthApi{
		"permission":   {Id: "GAAPI00000", Name: "認可ID発行API", Method: "GET", Endpoint: "/api/v1/permission"},
		"login":        {Id: "GAAPI00001", Name: "ログインAPI", Method: "POST", Endpoint: "/api/v1/login"},
		"get_token":    {Id: "GAAPI10000", Name: "アクセストークン要求API", Method: "POST", Endpoint: "/api/v1/token/create"},
		"check_token":  {Id: "GAAPI10001", Name: "トークン検証API", Method: "GET", Endpoint: "/api/v1/token/check"},
		"delete_token": {Id: "GAAPI10002", Name: "トークン削除API", Method: "DELETE", Endpoint: "/api/v1/token/delete"},
		"create_user":  {Id: "GAAPI20000", Name: "ユーザ登録API", Method: "POST", Endpoint: "/api/v1/user/create"},
		"get_user":     {Id: "GAAPI20001", Name: "ユーザ情報取得API", Method: "GET", Endpoint: "/api/v1/user/inquiry"},
	},
}

var ERROR_MESSAGE = map[string]string{
	"W0001": "リクエストメソッドが不適切です",
	"W0002": "リクエストパラメータが不適切です",
	"W0003": "認証エラーが発生しました",
	"W0005": "ユーザー名またはパスワードが違います",

	"E0001": "予期せぬエラーが発生しました",
}
