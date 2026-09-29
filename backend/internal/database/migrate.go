package database

import (
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"

	"github.com/fazasuny/erp-system/migrations"
)

// schemaMigrationsTable tracks which migration files have been applied. It is
// created by the migrator itself, not by a numbered SQL file, so existing
// databases that already have every table also work: the first boot records
// what it finds rather than failing.
const schemaMigrationsTable = "schema_migrations"

// Migrate applies every embedded migration that has not been recorded yet, in
// filename order, one transaction per file.
//
// Running this on boot is what closes the "fresh database has no tables" gap:
// previously migrations were manual (`psql -f`) and nothing in the binary
// executed them, so a new database started empty and the first query against
// `roles` failed with SQLSTATE 42P01.
//
// It is safe to call repeatedly. A migration that fails leaves nothing behind
// (its transaction rolls back) and is not recorded, so the server refuses to
// start instead of serving requests against a half-built schema.
func Migrate(db *gorm.DB) error {
	if err := ensureSchemaMigrationsTable(db); err != nil {
		return err
	}

	applied, err := appliedMigrations(db)
	if err != nil {
		return err
	}

	names, err := migrationNames()
	if err != nil {
		return err
	}

	// Databases migrated by hand before this migrator existed have all their
	// tables but no ledger. Running the files again would fail on "relation
	// already exists", so record them as already applied instead.
	if len(applied) == 0 && schemaAlreadyPresent(db) {
		log.Printf("existing schema detected without migration history; recording %d migrations as already applied", len(names))

		if err := db.Transaction(func(tx *gorm.DB) error {
			for _, name := range names {
				if err := tx.Exec(
					fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", schemaMigrationsTable),
					name,
				).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}

		return nil
	}

	pending := 0
	for _, name := range names {
		if _, ok := applied[name]; ok {
			continue
		}

		content, err := migrations.FS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		// One transaction per file: an early statement succeeding must not
		// leave a partially applied migration behind when a later one fails.
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(string(content)).Error; err != nil {
				return fmt.Errorf("apply migration %s: %w", name, err)
			}

			return tx.Exec(
				fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", schemaMigrationsTable),
				name,
			).Error
		}); err != nil {
			return err
		}

		log.Printf("migration applied: %s", name)
		pending++
	}

	if pending == 0 {
		log.Printf("database schema is up to date (%d migrations)", len(names))
	}

	return nil
}

// schemaAlreadyPresent reports whether the oldest migration's table exists,
// which is the signal that this database was built outside the migrator.
func schemaAlreadyPresent(db *gorm.DB) bool {
	var exists bool

	err := db.Raw(
		`SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'roles')`,
	).Scan(&exists).Error
	if err != nil {
		return false
	}

	return exists
}

func ensureSchemaMigrationsTable(db *gorm.DB) error {
	statement := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`, schemaMigrationsTable)

	return db.Exec(statement).Error
}

func appliedMigrations(db *gorm.DB) (map[string]struct{}, error) {
	var names []string
	if err := db.Raw(fmt.Sprintf("SELECT name FROM %s", schemaMigrationsTable)).Scan(&names).Error; err != nil {
		return nil, err
	}

	applied := make(map[string]struct{}, len(names))
	for _, name := range names {
		applied[name] = struct{}{}
	}

	return applied, nil
}

// migrationNames returns the embedded .sql files sorted by filename, which is
// the order the numeric prefixes encode.
func migrationNames() ([]string, error) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("list embedded migrations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		names = append(names, entry.Name())
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no migrations found in embedded filesystem")
	}

	sort.Strings(names)

	return names, nil
}