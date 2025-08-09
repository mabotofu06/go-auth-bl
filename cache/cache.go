package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
)

// cache はシングルトンでキャッシュインスタンスを保持
// アクセッサメソッドを通してアクセスする
var cache *ristretto.Cache[string, any]
var initNum = 0

func Init() error {
	if initNum > 0 {
		fmt.Println("キャッシュは初期化済です.")
		return nil
	}
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e7,     // 約1000万個のカウンター
		MaxCost:     1 << 30, // 1GBのキャッシュ
		BufferItems: 64,      // バッファのサイズ(基本64で十分)
	})
	if err != nil {
		initNum = 0
		return err
	}
	initNum++
	fmt.Println("キャッシュが初期化されました.")
	return nil
}

func checkCacheInitialized() error {
	if initNum <= 0 || cache == nil {
		return errors.New("キャッシュが初期化されていません")
	}
	return nil
}

// Setはキャッシュに値を設定します
func SetCache[T any](key string, value T, cost int64, ttl time.Duration) error {
	if err := checkCacheInitialized(); err != nil {
		return err
	}
	fmt.Printf("キャッシュに設定しました: %s\n", key)
	cache.SetWithTTL(key, value, cost, ttl)
	cache.Wait() // キャッシュの設定が完了するまで待つ
	return nil
}

// Getはキャッシュから値を取得します
func GetCache[T any](key string, deleteFlag bool) (T, bool) {
	var zero T
	if err := checkCacheInitialized(); err != nil {
		fmt.Println("キャッシュが初期化されていません:", err)
		return zero, false
	}

	value, found := cache.Get(key)
	if !found {
		return zero, false
	}

	v, ok := value.(T)
	if !ok {
		return zero, false
	}
	if deleteFlag {
		DeleteCache(key) // deleteFlagがtrueならキャッシュから削除
	}
	return v, true
}

// Delはキャッシュから値を削除します
func DeleteCache(key string) error {
	if err := checkCacheInitialized(); err != nil {
		return err
	}
	cache.Del(key)
	fmt.Printf("キャッシュから削除しました: %s\n", key)
	return nil
}
