package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SwaggerUIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("server/files/index.html")
		if err != nil {
			http.Error(w, "Swagger UI not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func SwaggerSpecHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("server/files/ccc_fixed.yaml")
		if err != nil {
			http.Error(w, "OpenAPI spec not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func DocumentationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("server/files/README.md")
		if err != nil {
			http.Error(w, "Documentation not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func ApiV2Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("server/files/ccc_fixed.yaml")
		if err != nil {
			http.Error(w, "OpenAPI spec not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func StaticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestedPath string
		if strings.HasPrefix(r.URL.Path, "/api/v2/static/") {
			requestedPath = strings.TrimPrefix(r.URL.Path, "/api/v2/static/")
		} else if strings.HasPrefix(r.URL.Path, "/api/v2/static") {
			requestedPath = strings.TrimPrefix(r.URL.Path, "/api/v2/static")
		} else {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		filename := strings.TrimSpace(requestedPath)
		if filename == "" || filename == "/" || filename == "@" {
			filename = "class_diagram.png"
		}
		filename = strings.TrimPrefix(filename, "@")
		filename = filepath.Base(filename)

		primary := filepath.Join("server", "files", "static", "png", filename)
		data, err := os.ReadFile(primary)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		ctype := http.DetectContentType(data)
		w.Header().Set("Content-Type", ctype)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func ReservedStaticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestedPath string
		if strings.HasPrefix(r.URL.Path, "/api/v2/reserved/") {
			requestedPath = strings.TrimPrefix(r.URL.Path, "/api/v2/reserved/")
		} else if strings.HasPrefix(r.URL.Path, "/api/v2/reserved") {
			requestedPath = strings.TrimPrefix(r.URL.Path, "/api/v2/reserved")
		} else {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		filename := strings.TrimSpace(requestedPath)
		if filename == "" || filename == "/" || filename == "@" {
			filename = "class_diagram.png"
		}
		filename = strings.TrimPrefix(filename, "@")
		filename = filepath.Base(filename)

		primary := filepath.Join("server", "files", "static", "png", filename)
		data, err := os.ReadFile(primary)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		ctype := http.DetectContentType(data)
		w.Header().Set("Content-Type", ctype)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
