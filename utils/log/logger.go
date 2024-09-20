package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// InitLogger Initialize the logger instance
func InitLogger() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:     os.Stdout,
		NoColor: false,
		FormatLevel: func(i interface{}) string {
			if i == nil {
				return "[INFO]"
			}
			return "[" + i.(string) + "]"
		},
		FormatMessage: func(i interface{}) string {
			if i == nil {
				return ""
			}
			return i.(string)
		},
		FormatTimestamp: func(i interface{}) string {
			if i == nil {
				return ""
			}
			return "[" + i.(string) + "]"
		},
	}

	consoleWriter.FormatCaller = func(i interface{}) string {
		if i == nil {
			return ""
		}
		// You can format the caller info here
		return "[" + i.(string) + "]"
	}

	// Initialize the logger with caller information
	log.Logger = zerolog.New(consoleWriter).With().
		Timestamp().
		Caller().
		Logger()
}

// GetLogger returns the configured zerolog logger
func GetLogger() zerolog.Logger {
	return log.Logger
}
