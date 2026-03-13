package a

import "log/slog"

func f() {
	slog.Info("запуск сервера") // want "log message should contain only english characters"
	slog.Info("starting server")
}
