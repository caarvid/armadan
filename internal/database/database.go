package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Create(db string) (*sql.DB, *sql.DB, error) {
	dbPath := fmt.Sprintf("file:%s", db)

	reader, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=ro&_fk=ON", dbPath))
	if err != nil {
		return nil, nil, err
	}

	reader.SetMaxOpenConns(100)

	writer, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=rwc&_txlock=immediate&_timeout=5000&_journal=WAL&_sync=NORMAL&_fk=ON", dbPath))
	if err != nil {
		return nil, nil, err
	}

	writer.SetMaxOpenConns(1)

	return reader, writer, nil
}
