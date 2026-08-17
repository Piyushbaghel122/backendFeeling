package src

import (
	"net/http"
	"os"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Base directory for our static frontend
	baseDir := "../client/out"
	
	// Determine the actual path of the requested file
	path := filepath.Join(baseDir, r.URL.Path)

	// Check if the file exists and is not a directory
	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		// File does not exist, so we serve index.html (SPA fallback)
		http.ServeFile(w, r, filepath.Join(baseDir, "index.html"))
		return
	}

	// File exists, serve it (e.g. JS, CSS, Images from _next)
	http.ServeFile(w, r, path)
}
