package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarvid/armadan/internal/database"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	create     bool
	email      string
	password   string
	role       string
	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string
)

func createPool() (*pgxpool.Pool, error) {
	return database.CreatePool(
		context.TODO(),
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
	)
}

func createUser(db schema.Querier) {
	ctx := context.TODO()
	hash, err := utils.GenerateHash(password, nil)
	if err != nil {
		log.Fatal(err)
	}

	user, err := db.CreateUser(ctx, &schema.CreateUserParams{
		Email:    email,
		Password: hash.Encode(),
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.UpdateUserRole(ctx, &schema.UpdateUserRoleParams{
		ID:   user.ID,
		Role: schema.UsersRoleEnum(role),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User created!")
}

func init() {
	flag.BoolVar(&create, "create", false, "create user")

	flag.StringVar(&email, "email", "", "user email")
	flag.StringVar(&password, "password", "", "user password")
	flag.StringVar(&role, "role", "user", "user role")

	flag.StringVar(&dbHost, "db_host", os.Getenv("DB_HOST"), "db host, defaults to env.DB_HOST")
	flag.StringVar(&dbPort, "db_port", os.Getenv("DB_PORT"), "db port, defaults to env.DB_PORT")
	flag.StringVar(&dbName, "db_name", os.Getenv("DB_NAME"), "db name, defaults to env.DB_NAME")
	flag.StringVar(&dbUser, "db_user", os.Getenv("DB_USER"), "db user, defaults to env.DB_USER")
	flag.StringVar(&dbPassword, "db_password", os.Getenv("DB_PASSWORD"), "db password, defaults to env.DB_PASSWORD")

	flag.Parse()
}

func main() {
	pool, err := createPool()
	if err != nil {
		log.Fatal(err)
	}

	db := schema.New(pool)

	if create {
		createUser(db)
	}
}
