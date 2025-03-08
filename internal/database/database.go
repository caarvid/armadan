package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Create() (*sql.DB, *sql.DB, error) {
	dbPath := "/Users/calle/projects/armadan/db/armadan.sqlite"

	reader, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro&_txlock=immediate", dbPath))
	if err != nil {
		return nil, nil, err
	}

	reader.SetMaxOpenConns(100)

	writer, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=rwc&_txlock=immediate", dbPath))
	if err != nil {
		return nil, nil, err
	}

	writer.SetMaxOpenConns(1)

	return reader, writer, nil
}
