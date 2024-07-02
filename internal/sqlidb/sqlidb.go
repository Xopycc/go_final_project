package sqlidb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func exists(dbPath string) bool {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("app path:", appPath)

	dbFile := filepath.Join(appPath, dbPath)
	fmt.Println("DB path:", dbFile)
	var exists bool
	_, err = os.Stat(dbFile)

	if err == nil {
		exists = true
	}

	return exists
}

func Open(driver, dataSrc string) (*sql.DB, error) {
	db, err := sql.Open(driver, dataSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	if exists(dataSrc) {
		return db, nil
	}

	_, err = db.Exec("create table scheduler(id integer PRIMARY KEY, date text NOT NULL, title text NOT NULL, comment text, repeat VARCHAR(128));")
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return db, nil
}
