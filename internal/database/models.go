package database

import (
	"time"
)

type Song struct {
	ID           int       `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Artist       string    `json:"artist" db:"artist"`
	Album        string    `json:"album" db:"album"`
	Duration     int       `json:"duration" db:"duration"`
	Fingerprint  string    `json:"fingerprint" db:"fingerprint"`
	HashSegments []string  `json:"hash_segments" db:"-"`
	DateAdded    time.Time `json:"date_added" db:"date_added"`
}

type MatchResult struct {
	IsMatch     bool    `json:"is_match"`
	Confidence  float64 `json:"confidence"`
	MatchOffset int     `json:"match_offset"`
	Song        *Song   `json:"song,omitempty"`
	TimeInSong  float64 `json:"time_in_song"`
}

type AudioFingerprint struct {
	TrackID      string   `json:"track_id"`
	TrackName    string   `json:"track_name"`
	Artist       string   `json:"artist"`
	Fingerprint  string   `json:"fingerprint"`
	HashSegments []string `json:"hash_segments"`
}

type RecordingStatus struct {
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Message   string `json:"message"`
	Countdown int    `json:"countdown,omitempty"`
}
