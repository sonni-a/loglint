package a

import (
	"log"
	"log/slog"
)

type Logger struct{}

func (Logger) Infof(format string, args ...interface{}) {}

func testSlog() {
	slog.Info("Starting server")      // want `log message must start with lowercase letter`
	slog.Info("запуск сервера")       // want `log message must be in English`
	slog.Info("server started!!!")    // want `log message contains special characters`
	slog.Info("user password: 12345") // want `log message may contain sensitive data`
	slog.Warn("Warning!")             // want `log message contains special characters` `log message must start with lowercase letter`
	slog.Error("Ошибка подключения")  // want `log message must be in English` `log message must start with lowercase letter`
	slog.Info("version 1.2.3")
	slog.Info("https://example.com/a")
	slog.Info("pass-word")
	slog.Info("pass!word")                     // want `log message contains special characters`
	slog.Info("done 🚀...") // want `log message must be in English` `log message contains special characters`
	Logger{}.Infof("User %s logged in", "bob") // want `log message must start with lowercase letter`
	log.Fatalf("error %s!!!", "x")             // want `log message contains special characters`
}
