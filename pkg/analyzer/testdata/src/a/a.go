package a

import (
	"log/slog"
)

func testSlog() {
	slog.Info("Starting server")      // want `log message must start with lowercase letter`
	slog.Info("запуск сервера")       // want `log message must be in English`
	slog.Info("server started!!!")    // want `log message contains special characters`
	slog.Info("user password: 12345") // want `log message may contain sensitive data`
	slog.Warn("Warning!")             // want `log message contains special characters`
	slog.Error("Ошибка подключения")  // want `log message must be in English`
}
