package dbutils

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	"github.com/joshdcuneo/go-ui/database"
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

func MigrateDatabase(db *sql.DB, callback ExecutionCallback) error {
	return executeGlob(db, "migrations/*.sql", callback)

}

func DropDatabase(dbPath string) error {
	err := os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func SeedDatabase(db *sql.DB, callback ExecutionCallback) error {
	return executeGlob(db, "seeds/*.sql", callback)
}

func executeGlob(db *sql.DB, glob string, callback ExecutionCallback) error {
	files, err := fs.Glob(database.EFS, glob)
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := fs.ReadFile(database.EFS, file)
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
