// internal/handlers/api.go
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

// RecordAudio starts the background recording process and immediately responds.
func (h *Handler) RecordAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Kick off the recording pipeline asynchronously
	go func() {
		h.recordingProcess()
	}()

	response := map[string]interface{}{
		"status":  "started",
		"message": "Recording started successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
}

// recordingProcess performs an early 3s guess and then the full-duration recording and matching.
func (h *Handler) recordingProcess() {
	// ---- EARLY 3s GUESS ----
	// Attempt a short capture, fingerprint, and best-match to push a quick song name to the UI.
	previewFile := fmt.Sprintf("%s/preview_%d.mp4", h.config.TempDir, time.Now().UnixNano())
	if err := audio.RecordScreenWithAudio(previewFile, 3); err == nil {
		if fp, err := audio.ExtractAudioFingerprint(previewFile); err == nil {
			if best, err := matching.FindBestMatch(h.db, fp); err == nil && best != nil && best.Song != nil && best.Song.Title != "" {
				h.broadcastWebSocketMessage("early_guess", map[string]string{"name": best.Song.Title})
			} else {
				h.broadcastWebSocketMessage("early_guess", map[string]string{"name": "Unknown"})
			}
		}
	}
	_ = os.Remove(previewFile)

	// ---- FULL RECORDING FLOW ----
	duration := h.config.RecordingDuration
	videoFile := fmt.Sprintf("%s/recording_%d.mp4", h.config.TempDir, time.Now().UnixNano())

	// Optional status broadcast (the UI already handles 'recording_status')
	h.broadcastStatus(database.RecordingStatus{
		Status:   "recording",
		Progress: 0,
		Message:  "Recording started...",
	})

	if err := audio.RecordScreenWithAudio(videoFile, duration); err != nil {
		h.broadcastStatus(database.RecordingStatus{
			Status:  "error",
			Message: fmt.Sprintf("Recording failed: %v", err),
		})
		return
	}

	fp, err := audio.ExtractAudioFingerprint(videoFile)
	if err != nil {
		h.broadcastStatus(database.RecordingStatus{
			Status:  "error",
			Message: fmt.Sprintf("Failed to process audio: %v", err),
		})
		_ = os.Remove(videoFile)
		return
	}

	result, err := matching.FindBestMatch(h.db, fp)
	if err != nil {
		h.broadcastStatus(database.RecordingStatus{
			Status:  "error",
			Message: fmt.Sprintf("Matching failed: %v", err),
		})
		_ = os.Remove(videoFile)
		return
	}

	h.broadcastResult(result)
	_ = os.Remove(videoFile)
}

// IdentifySong remains as a stub endpoint; the app uses /api/record + WebSocket for live flow.
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
	_ = json.NewEncoder(w).Encode(response)
}

// GetSongs streams all songs as JSON.
func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := h.db.GetAllSongs()
	if err != nil {
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(songs)
}

// AddSong ingests a song by artist/title/album using yt-dlp + fingerprinting pipeline.
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

	if err := audio.AddSongToDatabase(h.db, req.Artist, req.Title, req.Album); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add song: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Added %s - %s", req.Artist, req.Title),
	})
}

// SearchSongs performs a case-insensitive substring search on title/artist.
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
	_ = json.NewEncoder(w).Encode(songs)
}

// broadcastStatus forwards recording status to the frontend via WebSocket.
func (h *Handler) broadcastStatus(status database.RecordingStatus) {
	h.broadcastWebSocketMessage("recording_status", status)
}

// broadcastResult forwards the final match result to the frontend via WebSocket.
func (h *Handler) broadcastResult(result *database.MatchResult) {
	h.broadcastWebSocketMessage("result", result)
}
