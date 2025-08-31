package utils

import (
	"fmt"
	gelf "github.com/snovichkov/zap-gelf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"notification-service-api/internal/shared/rpc"
	"os"
	"strings"
)

var logger *zap.Logger
var level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

func parseLevel(s string) (zapcore.Level, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return zapcore.DebugLevel, true
	case "info", "":
		return zapcore.InfoLevel, true
	case "warn", "warning":
		return zapcore.WarnLevel, true
	case "error":
		return zapcore.ErrorLevel, true
	case "dpanic":
		return zapcore.DPanicLevel, true
	case "panic":
		return zapcore.PanicLevel, true
	case "fatal":
		return zapcore.FatalLevel, true
	default:
		return zapcore.InfoLevel, false
	}
}

func InitLogs() *zap.Logger {
	if lv, ok := parseLevel(os.Getenv("LOG_LEVEL")); ok {
		level.SetLevel(lv)
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout",
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		level,
	)

	consoleCore = rpc.WrapWithEndpointPrefix(consoleCore)

	var graylogEnabled bool
	var graylogAddr = strings.TrimSpace(os.Getenv("GRAYLOG_HOST"))

	cores := []zapcore.Core{consoleCore}
	if graylogAddr != "" {
		serviceName := "notification" + strings.ToLower(os.Getenv("GRAYLOG_SERVICE_POSTFIX"))

		graylogCore, err := gelf.NewCore(
			gelf.Addr(graylogAddr),
			gelf.Host(serviceName),
			gelf.LevelAtomic(level),
		)

		if err != nil {
			log.Printf("Failed to create Graylog core: %v", err)
		} else {
			graylogCore = rpc.WrapWithEndpointPrefix(graylogCore)

			cores = append(cores, graylogCore)
			graylogEnabled = true
		}
	}

	logger = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	host, _ := os.Hostname()

	logger = logger.With(
		zap.String("env", os.Getenv("SERVICE_ENV")),
		zap.String("gin_mode", os.Getenv("GIN_MODE")),
		zap.String("container_id", host),
	)

	if graylogEnabled && (os.Getenv("SERVICE_ENV") == "production" || os.Getenv("GRAYLOG_SMOKE_TEST") == "true") {
		logger.Info("logging initialized",
			zap.String("sink", "stdout+graylog"),
			zap.String("graylog_addr", graylogAddr),
			zap.String("level", level.Level().String()),
		)

		logger.Info("GRAYLOG_SMOKE_TEST", zap.String("marker", "zap-gelf"))
	} else {
		logger.Info("logging initialized", zap.String("sink", "stdout"), zap.String("level", level.Level().String()))
	}

	go func(l *zap.Logger) {
		if err := l.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}(logger)

	return logger
}

func SetLogLevel(newLevel string) error {
	lv, ok := parseLevel(newLevel)
	if !ok {
		return fmt.Errorf("unknown log level: %s", newLevel)
	}
	level.SetLevel(lv)
	logger.Info("Log level changed", zap.String("level", lv.String()))
	return nil
}

func GetLogger() *zap.Logger {
	return logger
}
