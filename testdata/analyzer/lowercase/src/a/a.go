package a

import "log/slog"

func f() {
	slog.Info("Starting server") // want "log message should start with lowercase letter"
	slog.Info("starting server")
}
