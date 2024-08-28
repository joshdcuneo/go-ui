package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joshdcuneo/go-ui/internal/dbutils"
)

var (
	migrationsDir string
	seedsDir      string
	seed          bool
	dbPath        string
)

func main() {
	flag.StringVar(&migrationsDir, "migrations-dir", "migrations", "Directory containing migration files")
	flag.StringVar(&seedsDir, "seeds-dir", "seeds", "Directory containing seed files")
	flag.BoolVar(&seed, "seed", false, "Run the seeders in addition to the action")
	flag.StringVar(&dbPath, "db", "db.sqlite", "Path to the SQLite database file")
	flag.Parse()

	if err := dbutils.EnsureDatabase(dbPath); err != nil {
		log.Fatal(err)
	}

	if flag.NArg() != 1 {
		log.Fatal("Usage: go run cmd/migrate [flags] <action>")
	}

	action := flag.Arg(0)

	logExecution := func(filename string) {
		fmt.Printf("\t%s\n", filename)
	}

	switch action {
	case "up":
		db, err := dbutils.NewDB(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		fmt.Println("Running migrations...")
		if err := dbutils.RunMigrations(db, migrationsDir, logExecution); err != nil {
			log.Fatal(err)
		}
		if seed {
			fmt.Println("Running seeds...")
			if err := dbutils.SeedDatabase(db, seedsDir, logExecution); err != nil {
				log.Fatal(err)
			}
		}
	case "drop":
		fmt.Println("Dropping database...")
		if err := dbutils.DropDatabase(dbPath); err != nil {
			log.Fatal(err)
		}
	case "fresh":
		fmt.Println("Dropping database...")
		if err := dbutils.DropDatabase(dbPath); err != nil {
			log.Fatal(err)
		}
		db, err := dbutils.NewDB(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		fmt.Println("Running migrations...")
		if err := dbutils.RunMigrations(db, migrationsDir, logExecution); err != nil {
			log.Fatal(err)
		}
		if seed {
			fmt.Println("Running seeds...")
			if err := dbutils.SeedDatabase(db, seedsDir, logExecution); err != nil {
				log.Fatal(err)
			}
		}
	case "seed":
		db, err := dbutils.NewDB(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		fmt.Println("Running seeds...")
		if err := dbutils.SeedDatabase(db, seedsDir, logExecution); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown action: %s", action)
	}

	fmt.Println("Done!")
}
