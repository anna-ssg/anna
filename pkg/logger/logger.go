package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
)

// Logger is a small wrapper around zerolog.Logger that exposes a few
// convenience methods with names compatible with the standard library's
// log.Logger used throughout the codebase (Printf, Println, Fatal, Fatalf).
type Logger struct {
	l zerolog.Logger
}

// New creates a new Logger writing to the provided io.Writer.
// If `out` is a terminal (TTY) a human-friendly, colorized console writer is used.
// Otherwise structured JSON logs are used.
func New(out io.Writer) *Logger {
	// If we're writing to a file that is a terminal, use the ConsoleWriter
	// which emits colorized, pretty output suitable for local development.
	if f, ok := out.(*os.File); ok && isatty.IsTerminal(f.Fd()) {
		console := zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: time.Kitchen, // e.g. "3:04PM"
			NoColor:    false,
		}
		z := zerolog.New(console).With().Timestamp().Logger()
		return &Logger{l: z}
	}

	// Fallback to structured JSON output (good for non-interactive logs).
	z := zerolog.New(out).With().Timestamp().Logger()
	return &Logger{l: z}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.l.Info().Msgf(format, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.l.Info().Msg(fmt.Sprint(v...))
}

func (l *Logger) Print(v ...interface{}) {
	l.l.Info().Msg(fmt.Sprint(v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.l.Fatal().Msg(fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.l.Fatal().Msgf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.l.Error().Msgf(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.l.Error().Msg(fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.l.Info().Msgf(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.l.Info().Msg(fmt.Sprint(v...))
}
