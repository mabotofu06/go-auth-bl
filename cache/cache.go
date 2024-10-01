package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

func SetupCache() (*ristretto.Cache[string, any], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e7,     // 約1000万個のカウンター
		MaxCost:     1 << 30, // 1GBのキャッシュ
		BufferItems: 64,      // バッファのサイズ
	})

	if err != nil {
		return nil, err
	}

	return cache, nil
}

// Setはキャッシュに値を設定します
func SetCache(cache *ristretto.Cache[string, any], key string, value any, cost int64, ttl time.Duration) {
	cache.SetWithTTL(key, value, cost, ttl)
	time.Sleep(10 * time.Millisecond) // 確実に設定されるまで待つ
}

// Getはキャッシュから値を取得します
func GetCache(cache *ristretto.Cache[string, any], key string) (any, bool) {
	return cache.Get(key)
}

// Delはキャッシュから値を削除します
func DeleteCache(cache *ristretto.Cache[string, any], key string) {
	cache.Del(key)
}
