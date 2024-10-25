/*
This Go code implements a function that automatically searches for a migration folder across the entire project and runs database migrations in a separate goroutine.
The process is broken down into the following steps:

1. **findMigrationFolder**: This function recursively searches the project directory for the folder containing migration files.
	It uses the `filepath.Walk` function to traverse the file tree.
	When the folder is found, it returns its path; otherwise, it returns an error if the folder is not found.

2. **getData**: This function initiates the migration process in a separate goroutine.
	It retrieves database connection parameters from environment variables (`DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`), then calls `findMigrationFolder` to locate the folder by name.
	Once the folder is found, the full migration path is formed, and migrations are run via the `runMigrations` function.

3. **runMigrations**: This function initializes and applies migrations using the `github.com/golang-migrate/migrate` package.
	It connects to the database using the provided connection string and applies the migrations.
	If there are no pending migrations (`migrate.ErrNoChange`), it simply returns without error.
*/

package migrate

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func findMigrationFolder(folderName string) (string, error) {
	log.Printf("Searching for migration folder: %s", folderName)
	var folderPath string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error during file walk: %v", err)
			return err
		}
		if info.IsDir() && strings.Contains(path, folderName) {
			folderPath = path
			log.Printf("Migration folder found at: %s", folderPath)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if folderPath == "" {
		return "", fmt.Errorf("folder %s not found", folderName)
	}
	return folderPath, nil
}

func Run(folderName string) {
	go func() {
		log.Println("Starting migration process...") // Log

		dbName := os.Getenv("DB_NAME")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")

		log.Printf("Connecting to database %s on host %s:%s", dbName, dbHost, dbPort) // Log

		dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

		migrationsPath, err := findMigrationFolder(folderName)
		if err != nil {
			log.Fatalf("Error finding migration folder: %v", err) // Fatal log if error
		}

		migrationsFullPath := fmt.Sprintf("file://%s", migrationsPath)
		log.Printf("Migrations path: %s", migrationsFullPath) // Log

		if err := runMigrations(dbURL, migrationsFullPath); err != nil {
			log.Fatalf("Error running migrations: %v", err)
		}

		log.Println("Migrations applied successfully!")
	}()
}

func runMigrations(dbURL, migrationsPath string) error {
	log.Println("Initializing migration...") // Log

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(migrate.ErrNoChange, err) {
		return fmt.Errorf("failed to apply migration: %w", err)
	}

	log.Println("Migration applied or already up to date") // Log success or no change
	return nil
}
