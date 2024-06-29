package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/caarvid/armadan/internal/handlers"
	m "github.com/caarvid/armadan/internal/middleware"
	"github.com/caarvid/armadan/internal/routes"
	"github.com/caarvid/armadan/internal/schema"
	"github.com/caarvid/armadan/internal/validation"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getDbString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable pool_max_conns=50",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
}

func createDatabasePool() (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(getDbString())

	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func initializeApp() *echo.Echo {
	app := echo.New()

	app.Validator = validation.NewValidator(validator.New(validator.WithRequiredStructEnabled()))
	app.Pre(middleware.RemoveTrailingSlash())
	app.Static("/public", "web/static")
	app.Use(m.DefaultContext)

	return app
}

func start() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	dbPool, err := createDatabasePool()

	if err != nil {
		return err
	}

	defer dbPool.Close()

	app := initializeApp()

	queries := schema.New(dbPool)
	handlers := handlers.Init(queries, dbPool)

	routes.Register(app, handlers, queries)

	return app.Start(":8080")
}

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}
