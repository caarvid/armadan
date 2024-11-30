package logger

import (
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func Create(logLevel zerolog.Level, isDev bool) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var output io.Writer = os.Stderr

	if isDev {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FieldsExclude: []string{
				"user_agent",
				"git_revision",
				"go_version",
			},
		}
	}

	buildInfo, _ := debug.ReadBuildInfo()

	return zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Str("go_version", buildInfo.GoVersion).
		Logger()
}
