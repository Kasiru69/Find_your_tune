package database

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func Initialize(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}

	if err := db.createTables(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) createTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS songs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        artist TEXT NOT NULL,
        album TEXT,
        duration INTEGER,
        fingerprint TEXT NOT NULL,
        hash_segments TEXT NOT NULL,
        date_added DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_artist ON songs(artist);
    CREATE INDEX IF NOT EXISTS idx_title ON songs(title);
    CREATE INDEX IF NOT EXISTS idx_fingerprint ON songs(fingerprint);
    `

	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) AddSong(song *Song) error {
	hashSegmentsJSON, err := json.Marshal(song.HashSegments)
	if err != nil {
		return err
	}

	query := `
    INSERT INTO songs (title, artist, album, duration, fingerprint, hash_segments)
    VALUES (?, ?, ?, ?, ?, ?)
    `

	result, err := db.conn.Exec(query, song.Title, song.Artist, song.Album,
		song.Duration, song.Fingerprint, string(hashSegmentsJSON))
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	song.ID = int(id)
	return nil
}

func (db *DB) GetAllSongs() ([]*Song, error) {
	query := `
    SELECT id, title, artist, album, duration, fingerprint, hash_segments, date_added
    FROM songs ORDER BY artist, title
    `

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*Song
	for rows.Next() {
		song := &Song{}
		var hashSegmentsJSON string
		var dateAdded string

		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.Album,
			&song.Duration, &song.Fingerprint, &hashSegmentsJSON, &dateAdded)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(hashSegmentsJSON), &song.HashSegments)
		songs = append(songs, song)
	}

	return songs, nil
}

func (db *DB) SearchSongs(query string) ([]*Song, error) {
	searchQuery := `
    SELECT id, title, artist, album, duration, fingerprint, hash_segments, date_added
    FROM songs 
    WHERE title LIKE ? OR artist LIKE ?
    ORDER BY artist, title
    `

	searchTerm := "%" + query + "%"
	rows, err := db.conn.Query(searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*Song
	for rows.Next() {
		song := &Song{}
		var hashSegmentsJSON string
		var dateAdded string

		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.Album,
			&song.Duration, &song.Fingerprint, &hashSegmentsJSON, &dateAdded)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(hashSegmentsJSON), &song.HashSegments)
		songs = append(songs, song)
	}

	return songs, nil
}

func (db *DB) GetSongCount() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM songs").Scan(&count)
	return count, err
}

func (db *DB) Close() error {
	return db.conn.Close()
}
