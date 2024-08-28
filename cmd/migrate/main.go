package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joshdcuneo/go-ui/internal/dbutils"
)

var (
	seed        bool
	databaseDSN string
)

func main() {
	flag.BoolVar(&seed, "seed", false, "Run the seeders in addition to the action")
	flag.StringVar(&databaseDSN, "database-dsn", "./db.sqlite", "Database DSN")
	flag.Parse()

	if err := dbutils.EnsureDatabase(databaseDSN); err != nil {
		log.Fatal(err)
	}

	if flag.NArg() != 1 {
		log.Fatal("Usage: migrate [flags] <action>")
	}

	action := flag.Arg(0)

	logExecution := func(filename string) {
		fmt.Printf("\t%s\n", filename)
	}

	switch action {
	case "up":
		db, err := dbutils.NewDB(databaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		fmt.Println("Running migrations...")
		if err := dbutils.MigrateDatabase(db, logExecution); err != nil {
			log.Fatal(err)
		}
		if seed {
			fmt.Println("Running seeds...")
			if err := dbutils.SeedDatabase(db, logExecution); err != nil {
				log.Fatal(err)
			}
		}
	case "drop":
		fmt.Println("Dropping database...")
		if err := dbutils.DropDatabase(databaseDSN); err != nil {
			log.Fatal(err)
		}
	case "seed":
		db, err := dbutils.NewDB(databaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		fmt.Println("Running seeds...")
		if err := dbutils.SeedDatabase(db, logExecution); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown action: %s", action)
	}

	fmt.Println("Done!")
}
