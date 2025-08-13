package database

import (
	"fmt"
)

func (db *DB) RunMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS songs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            artist TEXT NOT NULL,
            album TEXT,
            duration INTEGER,
            fingerprint TEXT NOT NULL,
            hash_segments TEXT NOT NULL,
            date_added DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE INDEX IF NOT EXISTS idx_artist ON songs(artist);`,
		`CREATE INDEX IF NOT EXISTS idx_title ON songs(title);`,
		`CREATE INDEX IF NOT EXISTS idx_fingerprint ON songs(fingerprint);`,
		`CREATE TABLE IF NOT EXISTS migration_history (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            migration_name TEXT NOT NULL,
            executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
	}

	for i, migration := range migrations {
		_, err := db.conn.Exec(migration)
		if err != nil {
			return fmt.Errorf("migration %d failed: %v", i+1, err)
		}
	}

	return nil
}

func (db *DB) GetMigrationHistory() ([]string, error) {
	query := `SELECT migration_name FROM migration_history ORDER BY executed_at`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		history = append(history, name)
	}

	return history, nil
}

func (db *DB) RecordMigration(name string) error {
	_, err := db.conn.Exec(`INSERT INTO migration_history (migration_name) VALUES (?)`, name)
	return err
}
