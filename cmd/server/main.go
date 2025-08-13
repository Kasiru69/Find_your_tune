package main

import (
	"log"
	"net/http"
	"os"

	"Shazam/config"
	"Shazam/internal/database"
	"Shazam/internal/handlers"
)

func main() {
	cfg := config.Load()

	db, err := database.Initialize(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	h := handlers.New(db, cfg)

	setupRoutes(h)

	os.MkdirAll("data/temp", 0755)

	log.Printf("🎵 Audio Recognition Server starting on http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

func setupRoutes(h *handlers.Handler) {
	http.HandleFunc("/", h.HomePage)
	http.HandleFunc("/database", h.DatabasePage)
	http.HandleFunc("/record", h.RecordPage)

	http.HandleFunc("/api/record", h.RecordAudio)
	http.HandleFunc("/api/identify", h.IdentifySong)
	http.HandleFunc("/api/songs", h.GetSongs)
	http.HandleFunc("/api/songs/add", h.AddSong)
	http.HandleFunc("/api/songs/search", h.SearchSongs)

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static/"))))

	http.HandleFunc("/ws", h.WebSocketHandler)
}
