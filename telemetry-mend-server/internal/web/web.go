package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var distFS embed.FS

func RegisterHandlers(r http.Handler) http.Handler {
	// Root of the embedded FS is 'dist'
	contentStatic, _ := fs.Sub(distFS, "dist")
	fileServer := http.FileServer(http.FS(contentStatic))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// If it's an API call, let the main router handle it
		if strings.HasPrefix(req.URL.Path, "/api") || req.URL.Path == "/health" {
			r.ServeHTTP(w, req)
			return
		}

		// Check if the file exists in the embedded FS
		path := strings.TrimPrefix(req.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		_, err := contentStatic.Open(path)
		if err != nil {
			// If file not found, serve index.html (SPA routing)
			req.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, req)
	})
}
