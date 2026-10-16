package c

import (
	"context"
	"log/slog"
)

type Server struct{ *slog.Logger }

type ZapClient struct{}

func (ZapClient) Info(msg string) {}

type Logger struct{}

func (Logger) InfoContext(msg, err string) {}
func (Logger) Log(msg, a, b string)        {}

type Entry struct{}

func (Entry) Info(msg string) {}

func cases() {
	var s Server
	s.Info("Starting server") // want `log message must start with lowercase letter`

	ZapClient{}.Info("HTTP status")
	Entry{}.Info("Starting server")

	Logger{}.InfoContext("Starting server", "ok") // want `log message must start with lowercase letter`
	Logger{}.InfoContext("ok", "BAD!!!")
	Logger{}.Log("BAD!!!", "ok", "ok") // want `log message contains special characters` `log message must start with lowercase letter`
	Logger{}.Log("ok", "ok", "BAD!!!")

	slog.InfoContext(context.Background(), "Starting server")                  // want `log message must start with lowercase letter`
	slog.Log(context.Background(), slog.LevelInfo, "Starting server")          // want `log message must start with lowercase letter`
}
