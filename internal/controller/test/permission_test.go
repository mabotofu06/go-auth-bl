package test

import (
	"go-auth-bl/internal/cache"
	"go-auth-bl/internal/controller"
	service_mock "go-auth-bl/internal/service/mock"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// クエリパラメータのテスト（正常、エラーで想定通りのレスポンスが返される）
func TestPermissionQueryParam(t *testing.T) {
	cache.Init()
	logger.Init("debug")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := service_mock.NewMockIClientInfoService(ctrl)
	mock.EXPECT().IsEnableClient(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	orig := controller.ClientInfoService
	controller.ClientInfoService = mock
	//終わった後は元のサービスに戻す
	t.Cleanup(func() { controller.ClientInfoService = orig })

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{"有効なクエリ(必須の値のみ)", "?response_type=code&client_id=123&redirect_uri=https://example.com/callback", http.StatusOK},
		{"有効なクエリ(scope付)", "?response_type=code&client_id=123&redirect_uri=https://example.com/callback&scope=email", http.StatusOK},
		{"有効なクエリ(state付)", "?response_type=code&client_id=123&redirect_uri=https://example.com/callback&state=xyz", http.StatusOK},
		{"無効なクエリ(response_typeがcode以外)", "?response_type=token", http.StatusBadRequest},
		{"クエリパラメータなし", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/permission"+tt.query, nil)
				w := httptest.NewRecorder()

				controller.GetPermission(w, req)

				res := w.Result()
				assert.Equal(t, tt.wantStatus, res.StatusCode)
			})
	}
}

// クライアント情報チェックのテスト（正常、エラーで想定通りのレスポンスが返される）
func TestPermissionClientCheck(t *testing.T) {
	cache.Init()
	logger.Init("debug")

	tests := []struct {
		name         string
		err          *ctm_err.CustomError
		expectStatus int
	}{
		{"クライアント情報チェックでエラーなし", nil, 200},
		{"クライアント情報チェックでサーバーエラー", ctm_err.InternalServerErr, 401},
		{"クライアント情報チェックで予期せぬエラー", ctm_err.UnexpectedServerErr, 401},
		{"クライアント情報チェックで予期せぬDBエラー", ctm_err.UnexpectedDBErr, 401},
		{"クライアント情報チェックでパラメータエラー", ctm_err.ParameterErr, 401},
		{"クライアント情報チェックで認証エラー", ctm_err.UnauthorizedErr, 401},
		{"クライアント情報チェックで存在しないエラー", ctm_err.NotFoundErr, 401},
		{"クライアント情報チェックで禁止エラー", ctm_err.ForbiddenErr, 401},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mock := service_mock.NewMockIClientInfoService(ctrl)
				mock.EXPECT().
					IsEnableClient(gomock.Any(), gomock.Any()).
					Return(tt.err).
					Times(1)

				orig := controller.ClientInfoService
				controller.ClientInfoService = mock
				t.Cleanup(func() { controller.ClientInfoService = orig })

				req := httptest.NewRequest(http.MethodGet, "/permission?response_type=code&client_id=123&redirect_uri=https://example.com/callback", nil)
				w := httptest.NewRecorder()

				controller.GetPermission(w, req)

				res := w.Result()
				assert.Equal(t, tt.expectStatus, res.StatusCode)
			})
	}
}
