package a

import "log/slog"

func f() {
	slog.Info("user password: secret") // want "log message may contain sensitive data"
	slog.Info("token: abc")            // want "log message may contain sensitive data"
	slog.Info("user authenticated successfully")
}
