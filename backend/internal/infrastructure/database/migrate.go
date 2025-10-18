package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes all SQL migration files in the migrations directory
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsPath string) error {
	log.Println("Running database migrations...")

	// Get all .up.sql files
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		log.Println("No migration files found")
		return nil
	}

	// Sort files to ensure they run in order
	sort.Strings(files)

	// Execute each migration file
	for _, file := range files {
		log.Printf("Executing migration: %s", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Printf("Migration completed: %s", filepath.Base(file))
	}

	log.Println("All migrations completed successfully")
	return nil
}
