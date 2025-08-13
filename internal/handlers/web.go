package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"Shazam/config"
	"Shazam/internal/database"
)

type Handler struct {
	db        *database.DB
	config    *config.Config
	templates map[string]*template.Template
}

func New(db *database.DB, cfg *config.Config) *Handler {
	h := &Handler{
		db:        db,
		config:    cfg,
		templates: make(map[string]*template.Template),
	}

	h.loadTemplates()
	return h
}

func (h *Handler) loadTemplates() {
	templateFiles := []string{"index.html", "database.html", "results.html"}

	for _, file := range templateFiles {
		tmpl := template.Must(template.ParseFiles(
			filepath.Join("web/templates", file),
		))
		h.templates[file] = tmpl
	}
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	count, _ := h.db.GetSongCount()

	data := struct {
		SongCount int
		Title     string
	}{
		SongCount: count,
		Title:     "Audio Recognition System",
	}

	h.templates["index.html"].Execute(w, data)
}

func (h *Handler) DatabasePage(w http.ResponseWriter, r *http.Request) {
	songs, err := h.db.GetAllSongs()
	if err != nil {
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	data := struct {
		Songs []database.Song
		Title string
	}{
		Songs: make([]database.Song, len(songs)),
		Title: "Song Database",
	}

	for i, song := range songs {
		data.Songs[i] = *song
	}

	h.templates["database.html"].Execute(w, data)
}

func (h *Handler) RecordPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Record & Identify",
	}

	h.templates["results.html"].Execute(w, data)
}
