package dbutils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type ExecutionCallback func(filename string)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsDir string, callback ExecutionCallback) error {
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("error executing migration %s: %v", file, err)
		}

		if callback != nil {
			callback(file)
		}
	}

	return nil
}

func DropDatabase(dbPath string) error {
	err := os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func SeedDatabase(db *sql.DB, seedsDir string, callback ExecutionCallback) error {
	files, err := filepath.Glob(filepath.Join(seedsDir, "*.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("error executing seed %s: %v", file, err)
		}

		if callback != nil {
			callback(file)
		}
	}

	return nil
}

func EnsureDatabase(dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("error creating database file: %v", err)
		}
		file.Close()
		fmt.Printf("Created database file: %s\n", dbPath)
	}
	return nil
}
