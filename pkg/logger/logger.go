package logger

import (
	"log/slog"
	"os"
)

func Start() {

	l := new(slog.HandlerOptions)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, l))

	slog.SetDefault(logger)
}
