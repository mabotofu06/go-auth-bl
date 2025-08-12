package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

var L *slog.Logger

var (
	infoLog  *slog.Logger
	errorLog *slog.Logger
)

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60) // JST (UTC+9)

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		// 既定の time フィールドを書き換え
		t := a.Value.Time().In(jst)
		//goの時間フォーマットは2006年1月2日15時4分5秒(固定値)のパターンを元に指定している
		a.Value = slog.StringValue(t.Format("2006/01/02 15:04:05.000"))
	}
	return a
}

// TODO:後々ファイル出力もできるように
func Init(level string) {
	lv := parseLevel(level)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lv, ReplaceAttr: replaceAttr})
	L = slog.New(handler)
	slog.SetDefault(L)
	slog.Info("Logger初期化完了")
}

func Print(level string, message string, args ...any) {
	switch level {
	case DEBUG:
		slog.Debug(fmt.Sprintf(message, args...))
	case INFO:
		slog.Info(fmt.Sprintf(message, args...))
	case WARN:
		slog.Warn(fmt.Sprintf(message, args...))
	case ERROR:
		slog.Error(fmt.Sprintf(message, args...))
	default:
		slog.Info(fmt.Sprintf(message, args...))
	}
}

func parseLevel(lv string) slog.Level {
	switch lv {
	case DEBUG:
		return slog.LevelDebug
	case INFO:
		return slog.LevelInfo
	case WARN:
		return slog.LevelWarn
	case ERROR:
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}
