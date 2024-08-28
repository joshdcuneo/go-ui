package models

import (
	"database/sql"
	"testing"

	"github.com/joshdcuneo/go-ui/internal/dbutils"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := dbutils.NewDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	dbutils.MigrateDatabase(db, "./../../migrations", nil)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
