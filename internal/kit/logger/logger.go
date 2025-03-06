package logger

import (
	"io"
	"log/slog"
)

const (
	sourceLocationKey = "logging.googleapis.com/sourceLocation"
)

func withGCPSeverity(groups []string, a slog.Attr) slog.Attr {
	switch {
	// TimeKey and format correspond to GCP convention by default
	// https://cloud.google.com/logging/docs/agent/logging/configuration#timestamp-processing
	case a.Key == slog.TimeKey && len(groups) == 0:
		return a
	case a.Key == slog.LevelKey && len(groups) == 0:
		logLevel, ok := a.Value.Any().(slog.Level)
		if !ok {
			return a
		}
		switch logLevel {
		case slog.LevelDebug:
			return slog.String("severity", "DEBUG")
		case slog.LevelInfo:
			return slog.String("severity", "INFO")
		case slog.LevelWarn:
			return slog.String("severity", "WARNING")
		case slog.LevelError:
			return slog.String("severity", "ERROR")
		default:
			return slog.String("severity", "DEFAULT")
		}
	case a.Key == slog.SourceKey && len(groups) == 0:
		source, ok := a.Value.Any().(*slog.Source)
		if !ok || source == nil {
			return a
		}

		return slog.Any(sourceLocationKey, source)
	case a.Key == slog.MessageKey && len(groups) == 0:

		return slog.String("message", a.Value.String())
	default:
		return a
	}
}

func NewLogger(w io.Writer, debug bool) *slog.Logger {
	lvl := slog.LevelInfo
	if debug {
		lvl = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       lvl,
		ReplaceAttr: withGCPSeverity,
	}))

	return logger
}
