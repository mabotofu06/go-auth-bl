package test

import (
	"go-auth-bl/internal/controller"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermission(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{"クエリパラメータなし", "", http.StatusBadRequest},
		{"有効なクエリ", "?user=alice", http.StatusOK}, // 期待するステータスに合わせて変更
		{"無効なクエリ", "?user=", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/permission"+tt.query, nil)
			w := httptest.NewRecorder()

			controller.GetPermission(w, req)

			res := w.Result()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
