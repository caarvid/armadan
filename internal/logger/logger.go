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

	var gitRevision string
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

	buildInfo, ok := debug.ReadBuildInfo()

	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
	}

	return zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Str("git_revision", gitRevision).
		Str("go_version", buildInfo.GoVersion).
		Logger()
}
