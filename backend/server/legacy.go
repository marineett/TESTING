package server

import (
	"net/http"
	"os"
	"path/filepath"
)

func LegacyArchiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		archivePath := filepath.Join("server", "files", "console2.zip")
		st, err := os.Stat(archivePath)
		if err != nil || st.IsDir() {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\"console2.zip\"")
		w.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
		http.ServeFile(w, r, archivePath)
	}
}
