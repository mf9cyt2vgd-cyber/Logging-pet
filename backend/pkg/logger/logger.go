package logger

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	logFile, err := os.OpenFile("../logs_pet/app.json.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		handler := slog.NewJSONHandler(os.Stdout, nil)
		return slog.New(handler)
	}

	// MultiWriter: пишем и в файл, и в stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	handler := slog.NewJSONHandler(multiWriter, nil)
	return slog.New(handler)
}
