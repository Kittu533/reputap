package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kittu533/reputap-be/internal/config"
	"github.com/kittu533/reputap-be/internal/database"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 1. Ensure migrations tracking table exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create schema_migrations table: %v", err)
	}

	// 2. Locate migrations directory (either ./migrations or ../migrations)
	migDir := "migrations"
	if _, err := os.Stat(migDir); os.IsNotExist(err) {
		migDir = filepath.Join("..", "migrations")
		if _, err := os.Stat(migDir); os.IsNotExist(err) {
			log.Fatalf("Migrations directory not found in ./migrations or ../migrations")
		}
	}

	files, err := os.ReadDir(migDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	fmt.Printf("📦 Found %d migration files in %s/\n", len(upFiles), migDir)

	for _, fname := range upFiles {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = $1`, fname).Scan(&count)
		if err != nil {
			log.Fatalf("Failed to check migration version %s: %v", fname, err)
		}

		if count > 0 {
			continue // Already applied
		}

		fmt.Printf("🚀 Applying migration: %s... ", fname)
		filePath := filepath.Join(migDir, fname)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", filePath, err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			log.Fatalf("\n❌ Failed applying %s: %v", fname, err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, fname); err != nil {
			tx.Rollback()
			log.Fatalf("\n❌ Failed recording migration %s: %v", fname, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("\n❌ Failed to commit migration %s: %v", fname, err)
		}

		fmt.Println("DONE")
	}

	fmt.Println("✅ All migrations applied successfully!")
}
