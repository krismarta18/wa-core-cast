package appserver

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:dashboard_out
var fsys embed.FS

func (s *Server) registerStaticRoutes() {
	// Get the static content from the embedded filesystem
	staticContent, err := fs.Sub(fsys, "dashboard_out")
	if err != nil {
		panic(err)
	}

	// Create a fileserver
	fileServer := http.FileServer(http.FS(staticContent))

	// Handler for static files and SPA fallback
	s.engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// If it's an API request, let Gin handle it normally (should have been caught by other routes)
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/health") {
			return
		}

		// 1. Check if exact file exists in embedded FS
		// Trim slashes to get a clean relative path
		cleanPath := strings.Trim(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}
		
		f, err := staticContent.Open(cleanPath)
		if err == nil {
			stat, _ := f.Stat()
			f.Close()
			// If it's a directory, we should still try the .html fallback 
			// because Next.js sometimes creates directories with the same name as pages
			if !stat.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// 2. Check if cleanPath + ".html" exists (for Next.js clean URLs)
		// This handles both /path and /path/ by looking for /path.html
		htmlPath := cleanPath + ".html"
		f, err = staticContent.Open(htmlPath)
		if err == nil {
			f.Close()
			// Update the request path so the file server can find the .html file
			c.Request.URL.Path = "/" + htmlPath
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// 3. SPA fallback: serve index.html
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
