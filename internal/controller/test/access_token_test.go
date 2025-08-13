package test

import (
	"bytes"
	"encoding/json"
	"go-auth-bl/internal/cache"
	"go-auth-bl/internal/controller"
	"go-auth-bl/internal/def"
	apiif "go-auth-bl/internal/dto/if"
	"go-auth-bl/internal/session"
	"go-auth-bl/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var def_req_body = controller.ReqAccessToken{
	Code:        "authcode-123",
	RedirectUri: "https://example.com/callback",
	State:       "xyz",
}

// アクセストークンの取得に成功した場合のテスト
func TestGetAccessToken_OK(t *testing.T) {
	cache.Init()
	logger.Init("debug")

	// 事前に認可コードのセッションをキャッシュへ投入
	code := "authcode-123"
	codeInfo := session.CodeInfo{
		ClientId:    "client-1",
		UserId:      "user-1",
		Scope:       "openid profile",
		RedirectUri: "https://example.com/callback",
		AccessToken: "access-xyz",
	}
	assert.NoError(t, cache.SetCache[session.CodeInfo](code, codeInfo, int64(1), 5*time.Minute))
	_, ok := cache.GetCache[session.CodeInfo](code, false)
	assert.True(t, ok, "code session should be cached")

	// リクエスト作成
	body, _ := json.Marshal(def_req_body)
	req := httptest.NewRequest(http.MethodPost, def.CONFIG.API["get_token"].Endpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 実行
	controller.GetAccessToken(w, req)

	// 検証
	res := w.Result()
	respBody := w.Body.Bytes()
	var got apiif.Response[controller.ResAccessToken]
	assert.NoError(t, json.Unmarshal(respBody, &got))

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Contains(t, res.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, codeInfo.AccessToken, got.Data.AccessToken)
	assert.Equal(t, codeInfo.UserId, got.Data.UserId)

	// 認可コードのセッションが削除されていることを確認
	_, ok = cache.GetCache[session.CodeInfo](code, false)
	assert.False(t, ok, "token has been deleted")

	// アクセストークンがキャッシュに保存されていること
	if tokenInfo, ok := cache.GetCache[session.TokenInfo](codeInfo.AccessToken, false); assert.True(t, ok, "token should be cached") {
		assert.Equal(t, codeInfo.ClientId, tokenInfo.ClientId)
		assert.Equal(t, codeInfo.UserId, tokenInfo.UserId)
		assert.Equal(t, codeInfo.Scope, tokenInfo.Scope)
	}
}

// アクセストークンのパラメータが不正な場合のテスト
func TestGetAccessToken_BadRequest_MissingParams(t *testing.T) {
	cache.Init()
	logger.Init("debug")

	cases := []struct {
		name    string
		reqBody controller.ReqAccessToken
	}{
		{
			name:    "コードが存在しない",
			reqBody: controller.ReqAccessToken{RedirectUri: "https://example.com/callback"},
		},
		{
			name:    "リダイレクトURIが存在しない",
			reqBody: controller.ReqAccessToken{Code: "authcode-123"},
		},
		{
			name:    "リクエストボディが空",
			reqBody: controller.ReqAccessToken{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// リクエスト作成
			body, _ := json.Marshal(tc.reqBody)
			req := httptest.NewRequest(http.MethodPost, def.CONFIG.API["get_token"].Endpoint, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 実行
			controller.GetAccessToken(w, req)

			// 検証
			res := w.Result()
			assert.Equal(t, 400, res.StatusCode)
		})
	}
}

// アクセストークンの認証エラー（コード未発行）のテスト
func TestGetAccessToken_AuthError_CodeNotFound(t *testing.T) {
	cache.Init()
	logger.Init("debug")

	// 存在しないコード
	body, _ := json.Marshal(def_req_body)
	req := httptest.NewRequest(http.MethodPost, def.CONFIG.API["get_token"].Endpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controller.GetAccessToken(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// アクセストークンのリダイレクトURIが不一致の場合のテスト
func TestGetAccessToken_AuthError_RedirectMismatch(t *testing.T) {
	cache.Init()
	logger.Init("debug")

	// 事前投入（redirect_uri を意図的に不一致にする）
	code := "authcode-123"
	codeInfo := session.CodeInfo{
		ClientId:    "client-2",
		UserId:      "user-2",
		Scope:       "email",
		RedirectUri: "https://example.com/callback-A",
		AccessToken: "access-abc",
	}
	assert.NoError(t, cache.SetCache[session.CodeInfo](code, codeInfo, int64(1), 5*time.Minute))

	body, _ := json.Marshal(def_req_body)
	req := httptest.NewRequest(http.MethodPost, def.CONFIG.API["get_token"].Endpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	controller.GetAccessToken(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}
