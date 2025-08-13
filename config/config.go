package config

import (
	"os"
)

type Config struct {
	Port                string
	DatabasePath        string
	TempDir             string
	SpotifyClientID     string
	SpotifyClientSecret string
	RecordingDuration   int
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		DatabasePath:        getEnv("DATABASE_PATH", "data/songs.db"),
		TempDir:             getEnv("TEMP_DIR", "data/temp"),
		SpotifyClientID:     getEnv("SPOTIFY_CLIENT_ID", "eec03041bad34931a01c2d8106bef880"),
		SpotifyClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", "66ea4b4480034839ae27ab41a9a20d1b"),
		RecordingDuration:   10,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
