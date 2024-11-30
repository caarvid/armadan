package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarvid/armadan/internal/database"
	"github.com/caarvid/armadan/internal/logger"
	"github.com/caarvid/armadan/internal/server"
	"github.com/caarvid/armadan/internal/service"
	"github.com/caarvid/armadan/pkg/assert"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/validation"
)

var (
	appEnv     string
	logLevel   string
	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string
	port       string
)

func run(
	ctx context.Context,
	args []string,
) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	dbPool, err := database.CreatePool(ctx, dbHost, dbPort, dbName, dbUser, dbPassword)
	if err != nil {
		return err
	}

	db := schema.New(dbPool)
	cache := cache.New(30*time.Minute, 15*time.Minute)

	srv := server.New(
		service.NewPostService(db, cache),
		service.NewWeekService(db, cache),
		service.NewUserService(db),
		service.NewPlayerService(db, dbPool),
		service.NewSessionService(db),
		service.NewCourseService(db, dbPool, cache),
		service.NewResultService(db, dbPool),
		validation.New(),
	)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: srv,
	}

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		log.Info().Msgf("server listening on %s", httpServer.Addr)

		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("error while listening")
			return err
		}

		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()

		log.Info().Msg("server is shutting down")

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("error while shutting down")
			return err
		}

		log.Info().Msg("server shutdown complete")
		return nil
	})

	return eg.Wait()
}

var logLevelMap = map[string]zerolog.Level{
	"DEBUG": zerolog.DebugLevel,
	"INFO":  zerolog.InfoLevel,
	"WARN":  zerolog.WarnLevel,
	"ERROR": zerolog.ErrorLevel,
	"FATAL": zerolog.FatalLevel,
	"OFF":   zerolog.Disabled,
}

func init() {
	flag.StringVar(&appEnv, "env", "development", "app environment")
	flag.StringVar(&logLevel, "log_level", "INFO", "log level")

	flag.StringVar(&dbHost, "db_host", os.Getenv("DB_HOST"), "db host, defaults to env.DB_HOST")
	flag.StringVar(&dbPort, "db_port", os.Getenv("DB_PORT"), "db port, defaults to env.DB_PORT")
	flag.StringVar(&dbName, "db_name", os.Getenv("DB_NAME"), "db name, defaults to env.DB_NAME")
	flag.StringVar(&dbUser, "db_user", os.Getenv("DB_USER"), "db user, defaults to env.DB_USER")
	flag.StringVar(&dbPassword, "db_password", os.Getenv("DB_PASSWORD"), "db password, defaults to env.DB_PASSWORD")
	flag.StringVar(&port, "port", os.Getenv("PORT"), "port, defaults to env.PORT")

	flag.Parse()
}

func main() {
	assert.OneOf(appEnv, []string{"development", "production", "test"}, "env must be one of [development, production, test]")
	assert.OneOf(logLevel, []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "OFF"}, "log_level must be one of [DEBUG, INFO, WARN, ERROR, FATAL, OFF]")

	ctx := context.Background()
	log.Logger = logger.Create(logLevelMap[logLevel], appEnv == "development")

	if err := run(ctx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("an unexpected error occurred")
	}
}
