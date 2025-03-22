package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/database"
	"github.com/caarvid/armadan/internal/database/schema"
	"github.com/caarvid/armadan/internal/utils"
)

var (
	create   bool
	email    string
	password string
	role     string
)

func createUser(writer schema.Querier) {
	ctx := context.TODO()
	hash, err := utils.GenerateHash(password, nil)
	if err != nil {
		log.Fatal(err)
	}

	user, err := writer.CreateUser(ctx, &schema.CreateUserParams{
		ID:       armadan.GetId(),
		Email:    email,
		Password: hash.Encode(),
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.UpdateUserRole(ctx, &schema.UpdateUserRoleParams{
		ID:       user.ID,
		UserRole: role,
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

	flag.Parse()
}

func main() {
	_, writeDB, err := database.Create()
	if err != nil {
		log.Fatal(err)
	}

	dbWriter := schema.New(writeDB)

	if create {
		createUser(dbWriter)
	}
}
