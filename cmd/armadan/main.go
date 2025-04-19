package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database"
	"github.com/caarvid/armadan/internal/logger"
	"github.com/caarvid/armadan/internal/server"
	"github.com/caarvid/armadan/internal/service"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/validation"
)

var (
	appEnv      string
	logLevel    string
	dbPath      string
	port        string
	logLevelMap = map[string]zerolog.Level{
		"DEBUG": zerolog.DebugLevel,
		"INFO":  zerolog.InfoLevel,
		"WARN":  zerolog.WarnLevel,
		"ERROR": zerolog.ErrorLevel,
		"FATAL": zerolog.FatalLevel,
		"OFF":   zerolog.Disabled,
	}
)

func run(
	ctx context.Context,
	getEnv func(string) string,
	_ []string,
) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	readDB, writeDB, err := database.Create(dbPath)
	if err != nil {
		return err
	}

	defer readDB.Close()
	defer writeDB.Close()

	awsConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return err
	}

	dbReader := schema.New(readDB)
	dbWriter := schema.New(writeDB)
	cache := cache.New(30*time.Minute, 15*time.Minute)
	emailOverride := getEnv("EMAIL_TARGET_OVERRIDE")

	if appEnv == "production" {
		emailOverride = ""
	}

	srv := server.New(
		service.NewPostService(dbReader, dbWriter, cache),
		service.NewWeekService(dbReader, dbWriter, cache),
		service.NewUserService(dbReader, dbWriter),
		service.NewPlayerService(dbReader, dbWriter, writeDB),
		service.NewSessionService(dbReader, dbWriter),
		service.NewCourseService(dbReader, dbWriter, writeDB, cache),
		service.NewResultService(dbReader, dbWriter, writeDB, cache),
		service.NewResetPasswordService(dbReader, dbWriter, writeDB),
		service.NewEmailService(
			armadan.Senders{
				ResetPassword: getEnv("RESET_PASSWORD_EMAIL_SENDER"),
			},
			awsConfig,
			emailOverride,
		),
		validation.New(),
	)

	httpServer := http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           srv,
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

func setDefaultStringValue(key string, f *string, allowedValues []string, fallback string) {
	if slices.Contains(allowedValues, *f) {
		return
	}

	fmt.Printf("%s must be one of [%s], falling back to %s", key, strings.Join(allowedValues, ", "), fallback)

	*f = fallback
}

func init() {
	flag.StringVar(&appEnv, "env", os.Getenv("APP_ENV"), "app environment, default to env.APP_ENV")
	flag.StringVar(&logLevel, "log_level", os.Getenv("LOG_LEVEL"), "log level, default to env.LOG_LEVEL")
	flag.StringVar(&dbPath, "db_path", os.Getenv("DB_PATH"), "path to sqlite db, defaults to env.DB_PATH")
	flag.StringVar(&port, "port", os.Getenv("PORT"), "port, defaults to env.PORT")

	flag.Parse()
}

func main() {
	setDefaultStringValue("env", &appEnv, []string{"development", "production", "test"}, "production")
	setDefaultStringValue("log level", &logLevel, []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "OFF"}, "ERROR")

	ctx := context.Background()
	log.Logger = logger.Create(logLevelMap[logLevel], appEnv == "development")

	if err := run(ctx, os.Getenv, os.Args); err != nil {
		log.Fatal().Err(err).Msg("an unexpected error occurred")
	}
}
