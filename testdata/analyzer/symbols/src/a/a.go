package a

import "log/slog"

func f() {
	slog.Info("server started!!!") // want "log message contains disallowed symbols"
	slog.Info("server started 🚀")  // want "log message contains disallowed symbols"
	slog.Info("server started")
}
