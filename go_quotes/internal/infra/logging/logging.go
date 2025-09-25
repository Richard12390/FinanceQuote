package logging

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(path string, level zapcore.Level, maxSizeMB, maxBackups, maxAgeDays int) (*zap.Logger, *lumberjack.Logger, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0o755)
	}
	rot := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   false,
	}
	enc := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(enc), zapcore.AddSync(rot), level)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enc),
		zapcore.AddSync(os.Stdout),
		level,
	)
	core := zapcore.NewTee(fileCore, consoleCore)

	log := zap.New(core, zap.AddCaller())
	return log, rot, nil
}

func StartDailyRotate(ctx context.Context, rot *lumberjack.Logger, loc *time.Location) {
	if rot == nil {
		return
	}
	if loc == nil {
		loc = time.Local
	}
	go func() {
		for {
			now := time.Now().In(loc)
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
			t := time.NewTimer(time.Until(next))
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				_ = rot.Rotate()
			}
		}
	}()
}
