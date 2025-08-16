package common_db

import (
	"context"
	"database/sql"
	"fmt"
	ctm_err "go-auth-bl/pkg/error"
	"go-auth-bl/pkg/logger"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// データベースに接続
func ConnectDB() (*sql.DB, *ctm_err.CustomError) {
	start := time.Now()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// 旧来の個別指定から組み立て
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			valueOr(os.Getenv("DB_SSLMODE"), "verify-full"),
		)
	}

	masked := maskPassword(dsn)
	logger.Print(logger.DEBUG, "DB接続開始 DSN=%s", masked)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Print(logger.ERROR, "sql.Open失敗: %v", err)
		return nil, ctm_err.UnexpectedDBErr
	}
	// コネクションプール調整 (Supabase Free 想定: 少なめ)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err = db.PingContext(ctx); err != nil {
		logger.Print(logger.ERROR, "Ping失敗: %v", err)
		return nil, ctm_err.UnexpectedDBErr
	}

	elapsed := time.Since(start)
	logger.Print(logger.DEBUG, "DB接続成功 elapsed=%s open=%d inuse=%d idle=%d",
		elapsed,
		db.Stats().OpenConnections,
		db.Stats().InUse,
		db.Stats().Idle,
	)
	logger.Print(logger.INFO, "データベース接続確立")
	return db, nil
}

func valueOr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func maskPassword(dsn string) string {
	// URL 形式か key=value 群かを判定
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		u, err := url.Parse(dsn)
		if err != nil || u.User == nil {
			return dsn
		}
		if _, ok := u.User.Password(); ok {
			u.User = url.UserPassword(u.User.Username(), "****")
		}
		return u.String()
	}
	parts := strings.Fields(dsn)
	for i, p := range parts {
		if strings.HasPrefix(p, "password=") {
			parts[i] = "password=****"
		}
	}
	return strings.Join(parts, " ")
}
