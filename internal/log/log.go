package log

import (
	"io"
	"log"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(io.Discard).With().Timestamp().Logger().Level(zerolog.InfoLevel)

func Human(w io.Writer) {
	logger = logger.Output(zerolog.ConsoleWriter{Out: w})
}

func Disable() {
	logger = logger.Level(zerolog.Disabled)
}

func Output(w io.Writer) {
	logger = logger.Output(w)
}

func SetDebug() {
	logger = logger.Level(zerolog.DebugLevel)
}

func SetInfo() {
	logger = logger.Level(zerolog.InfoLevel)
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Error() *zerolog.Event {
	return logger.Error().Caller(zerolog.CallerSkipFrameCount - 1)
}

func Warn() *zerolog.Event {
	return logger.Warn().Caller()
}

func Err(err error) *zerolog.Event {
	return logger.Err(err).Caller(zerolog.CallerSkipFrameCount - 1)
}

func Stdlib(name string) *log.Logger {
	return log.New(logger.With().Str("from", name).Logger(), "", 0)
}
