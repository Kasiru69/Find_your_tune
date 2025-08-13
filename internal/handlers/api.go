package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"Shazam/internal/audio"
	"Shazam/internal/database"
	"Shazam/internal/matching"
)

func (h *Handler) RecordAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	go func() {
		h.recordingProcess()
	}()

	response := map[string]interface{}{
		"status":  "started",
		"message": "Recording started successfully",
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) recordingProcess() {
	duration := h.config.RecordingDuration
	videoFile := fmt.Sprintf("%s/recording_%d.mp4", h.config.TempDir, time.Now().Unix())

	err := audio.RecordScreenWithAudio(videoFile, duration)
	if err != nil {
		h.broadcastStatus(database.RecordingStatus{
			Status:  "error",
			Message: fmt.Sprintf("Recording failed: %v", err),
		})
		return
	}

	fingerprint, err := audio.ExtractAudioFingerprint(videoFile)
	if err != nil {
		h.broadcastStatus(database.RecordingStatus{
			Status:  "error",
			Message: fmt.Sprintf("Failed to process audio: %v", err),
		})
		os.Remove(videoFile)
		return
	}

	result, err := matching.FindBestMatch(h.db, fingerprint)
	if err != nil {
		h.broadcastStatus(database.RecordingStatus{
			Status:  "error",
			Message: fmt.Sprintf("Matching failed: %v", err),
		})
		os.Remove(videoFile)
		return
	}

	h.broadcastResult(result)
	os.Remove(videoFile)
}

func (h *Handler) IdentifySong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"song": map[string]interface{}{
			"title":      "Shape of You",
			"artist":     "Ed Sheeran",
			"confidence": 67.3,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := h.db.GetAllSongs()
	if err != nil {
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Artist string `json:"artist"`
		Title  string `json:"title"`
		Album  string `json:"album"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Artist == "" || req.Title == "" {
		http.Error(w, "Artist and title are required", http.StatusBadRequest)
		return
	}

	err := audio.AddSongToDatabase(h.db, req.Artist, req.Title, req.Album)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add song: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Added %s - %s", req.Artist, req.Title),
	})
}

func (h *Handler) SearchSongs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	songs, err := h.db.SearchSongs(query)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *Handler) broadcastStatus(status database.RecordingStatus) {
}

func (h *Handler) broadcastResult(result *database.MatchResult) {
}
