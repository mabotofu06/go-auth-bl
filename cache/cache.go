package cache

import (
	"errors"
	"go-auth-bl/internal/session"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

var (
	cacheInstance *ristretto.Cache[string, any]
	once          sync.Once
	initError     error
)

// GetCacheInstance はシングルトンでキャッシュインスタンスを取得します
func GetCacheInstance() (*ristretto.Cache[string, any], error) {
	once.Do(func() {
		cache, err := ristretto.NewCache(&ristretto.Config[string, any]{
			NumCounters: 1e7,     // 約1000万個のカウンター
			MaxCost:     1 << 30, // 1GBのキャッシュ
			BufferItems: 64,      // バッファのサイズ
		})

		if err != nil {
			initError = err
			return
		}

		cacheInstance = cache
	})

	if initError != nil {
		return nil, initError
	}

	return cacheInstance, nil
}

// 以下は後方互換性のために残す（非推奨）
func SetupCache() (*ristretto.Cache[string, any], error) {
	return GetCacheInstance()
}

// Setはキャッシュに値を設定します
func SetCache(key string, value any, cost int64, ttl time.Duration) error {
	cache, err := GetCacheInstance()
	if err != nil {
		return err
	}

	cache.SetWithTTL(key, value, cost, ttl)
	time.Sleep(10 * time.Millisecond) // 確実に設定されるまで待つ
	return nil
}

// Getはキャッシュから値を取得します
func GetCache(key string) (any, bool) {
	cache, err := GetCacheInstance()
	if err != nil {
		return nil, false
	}

	return cache.Get(key)
}

// Delはキャッシュから値を削除します
func DeleteCache(key string) error {
	cache, err := GetCacheInstance()
	if err != nil {
		return err
	}

	cache.Del(key)
	return nil
}

// アクセストークンでセッション情報を取得
func GetTokenFromCache(accessToken string) (*session.TokenInfo, error) {
	value, found := GetCache(accessToken)
	if !found {
		return nil, errors.New("token not found")
	}

	tokenInfo, ok := value.(*session.TokenInfo)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return tokenInfo, nil
}
