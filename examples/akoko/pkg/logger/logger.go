package logger

import (
	"io"
	"log/slog"
	"os"
)

var AkokoLog *slog.Logger

func NewAkokoLogger(logfile *os.File) *slog.Logger {
	multiWritter := io.MultiWriter(os.Stderr, logfile)
	log := slog.New(slog.NewJSONHandler(multiWritter, &slog.HandlerOptions{
		AddSource: true,
	}))
	return log
}
