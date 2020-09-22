package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mikestefanello/otcscanner/config"
	"github.com/mikestefanello/otcscanner/handlers"
)

// NewRouter returns a new router
func NewRouter(cfg config.Config, h *handlers.HTTPHandler) *chi.Mux {
	r := chi.NewRouter()

	// Add basic auth
	if cfg.HTTP.Auth.User != "" && cfg.HTTP.Auth.Password != "" {
		auth := make(map[string]string)
		auth[cfg.HTTP.Auth.User] = cfg.HTTP.Auth.Password
		r.Use(middleware.BasicAuth("app", auth))
	}

	// Add routes
	r.Get("/", h.ScanForm)
	r.Post("/", h.ScanForm)
	r.Get("/database", h.DatabasePage)
	r.Post("/database/upload", h.DatabaseUpload)
	r.Post("/database/delete/all", h.DatabaseDeleteAll)
	r.Post("/database/delete/complete", h.DatabaseDeleteCompleted)
	r.Post("/database/download/all", h.DatabaseDownloadAll)
	r.Post("/database/download/completed", h.DatabaseDownloadCompleted)
	r.Post("/database/download/incomplete", h.DatabaseDownloadIncomplete)

	return r
}
