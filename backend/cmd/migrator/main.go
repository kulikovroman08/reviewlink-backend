package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const defaultMigrationsPath = "migrations"

func main() {
	log.Println("start migrator")

	_ = godotenv.Load()

	var migrationsPath string
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations directory (default: migrations)")
	flag.Parse()

	if migrationsPath == "" {
		migrationsPath = defaultMigrationsPath
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		log.Fatal("DB_URL (or DATABASE_URL) is empty")
	}

	absMigrationsPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Fatalf("failed to resolve migrations path: %v", err)
	}
	sourceURL := "file://" + absMigrationsPath

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("migrate source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("migrate db close error: %v", dbErr)
		}
	}()

	// Show state + fail fast on dirty
	if v, dirty, err := m.Version(); err == nil {
		log.Printf("current version=%d dirty=%v", v, dirty)
		if dirty {
			log.Fatalf("database is dirty at version %d; fix it and then run migrate force if needed", v)
		}
	} else if errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("current version=<nil> dirty=false (no migrations applied yet)")
	} else {
		log.Printf("could not read current version: %v", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}
		log.Fatalf("migration up failed: %v", err)
	}

	log.Println("migrations applied")
}
