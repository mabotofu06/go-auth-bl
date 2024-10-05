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
		"login":           {name: "認証", method: "POST", endpoint: "/api/login"},
		"access_token":    {name: "認可ID発行", method: "GET", endpoint: "/api/access_token"},
	},
}
