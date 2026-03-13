package a

import (
	"go.uber.org/zap"
	"log/slog"
)

func f() {
	slog.Info("Starting server")       // want "log message should start with lowercase letter"
	slog.Info("запуск сервера")        // want "log message should contain only english characters"
	slog.Info("server started!!!")     // want "log message contains disallowed symbols"
	slog.Info("user password: secret") // want "log message may contain sensitive data"
	slog.Info("starting server")

	zap.L().Info("Starting zap server") // want "log message should start with lowercase letter"
	zap.L().Info("zap server started")
}
