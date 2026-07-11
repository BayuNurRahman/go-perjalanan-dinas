package config

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is the global application logger.
var Logger *slog.Logger

// InitLogger initializes the global structured logger.
// Logs are written to both stdout and a rotating log file via Lumberjack.
func InitLogger() {
	// Ensure the logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		slog.Error("Gagal membuat direktori logs", "error", err)
	}

	// Configure Lumberjack for rotating file logs
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,   // Max megabytes before rotation
		MaxBackups: 5,    // Max number of old log files to retain
		MaxAge:     30,   // Max days to retain old log files
		Compress:   true, // Compress rotated log files
	}

	// Write to both stdout (text for readability) and log file (JSON for parsing)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Use JSON handler so logs in the file are structured and parseable
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: getLogLevel(),
	})

	Logger = slog.New(handler)
	slog.SetDefault(Logger)

	slog.Info("Logger berhasil diinisialisasi", "file", "logs/app.log")
}

// getLogLevel reads the LOG_LEVEL environment variable and returns the appropriate slog.Level.
// Supported values: "debug", "warn", "error". Defaults to "info".
func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
