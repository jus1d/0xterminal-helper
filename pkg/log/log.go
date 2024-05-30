package log

import (
	"fmt"
	"log/slog"
	"os"
	"terminal/internal/config"
	"terminal/pkg/log/prettyslog"
	"time"
)

// Init initialize a *slog.Logger instance for logging, without pretty formatting for production and development builds.
func Init(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = InitPretty()
	case config.EnvDevelopment:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProduction:
		logsDir := "./.logs"
		if _, err := os.Stat(logsDir); os.IsNotExist(err) {
			os.Mkdir(logsDir, 0755)
		}

		out, _ := os.OpenFile(time.Now().Format(fmt.Sprintf("%s/02-01-2006.log", logsDir)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		log = slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// InitPretty initialize a *slog.Logger instance for logging, with pretty formatting for local builds.
func InitPretty() *slog.Logger {
	opts := prettyslog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	prettyHandler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(prettyHandler)
}
