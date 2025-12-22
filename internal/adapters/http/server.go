package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
)

// ServerConfig contient la configuration du serveur
type ServerConfig struct {
	Port               int
	IsProd             bool
	EnableHSTS         bool
	BasicAuthConfig    BasicAuthConfig
	RateLimitPerMinute int
	MaxBodySize        int64
	StaticDir          string
	TemplatesDir       string
}

// Server représente le serveur HTTP
type Server struct {
	config   ServerConfig
	handlers *Handlers
	router   *chi.Mux
	server   *http.Server
}

// NewServer crée un nouveau serveur
func NewServer(config ServerConfig, handlers *Handlers) *Server {
	return &Server{
		config:   config,
		handlers: handlers,
	}
}

// setupRoutes configure les routes
func (s *Server) setupRoutes() {
	r := chi.NewRouter()

	// Rate limiter global
	rateLimiter := NewRateLimiter(s.config.RateLimitPerMinute)

	// Middlewares globaux
	globalMiddlewares := []Middleware{
		requestIDMiddleware(),
		recoverMiddleware(s.config.IsProd),
		loggingMiddleware(),
		securityHeadersMiddleware(s.config.EnableHSTS),
		maxBytesMiddleware(s.config.MaxBodySize),
		timeoutMiddleware(30 * time.Second),
	}

	// Basic Auth si configuré
	if s.config.BasicAuthConfig.Enabled {
		globalMiddlewares = append(globalMiddlewares, basicAuthMiddleware(s.config.BasicAuthConfig))
	}

	// Fichiers statiques (sans rate limiting pour les assets)
	fileServer := http.FileServer(http.Dir(s.config.StaticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Routes publiques avec rate limiting léger
	r.Group(func(r chi.Router) {
		r.Use(s.adaptMiddleware(Chain(
			func(w ResponseWriter, r *Request) error { return nil },
			append(globalMiddlewares, cacheControlMiddleware())...,
		)))

		r.Get("/", s.adaptHandler(s.handlers.HomeHandler, globalMiddlewares))
		r.Get("/planning", s.adaptHandler(s.handlers.PlanningHandler, globalMiddlewares))
		r.Get("/infos", s.adaptHandler(s.handlers.InfosHandler, globalMiddlewares))
		r.Get("/calendar.ics", s.adaptHandler(s.handlers.CalendarHandler, globalMiddlewares))
		r.Get("/health", s.adaptHandler(s.handlers.HealthHandler, globalMiddlewares))
		r.Get("/admin", s.adaptHandler(s.handlers.AdminHandler, globalMiddlewares))
		r.Get("/admin/export", s.adaptHandler(s.handlers.AdminExportHandler, globalMiddlewares))
		r.Get("/admin/delete", s.adaptHandler(s.handlers.AdminDeleteHandler, globalMiddlewares))
	})

	// Routes RSVP avec rate limiting strict
	r.Group(func(r chi.Router) {
		strictMiddlewares := append(globalMiddlewares, rateLimitMiddleware(rateLimiter))
		r.Get("/rsvp", s.adaptHandler(s.handlers.RSVPGetHandler, strictMiddlewares))
		r.Post("/rsvp", s.adaptHandler(s.handlers.RSVPPostHandler, strictMiddlewares))
	})

	// 404 handler
	r.NotFound(s.adapt(s.handlers.NotFoundHandler, globalMiddlewares))

	s.router = r
}

// adaptHandler adapte un HandlerFunc custom vers http.Handler
func (s *Server) adaptHandler(h HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	handler := Chain(h, middlewares...)
	return s.adapt(handler, []Middleware{})
}

// adapt convertit un HandlerFunc en http.HandlerFunc
func (s *Server) adapt(h HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	handler := Chain(h, middlewares...)

	return func(w http.ResponseWriter, r *http.Request) {
		// Wrapper le ResponseWriter
		rw := &responseWriter{ResponseWriter: w, status: 0}

		// Wrapper la Request
		req := &Request{Request: r}

		// Exécuter le handler
		if err := handler(rw, req); err != nil {
			log.Printf("Handler error (request_id=%s): %v", req.RequestID(), err)

			if rw.Status() == 0 {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			}
		}
	}
}

// adaptMiddleware adapte un middleware pour chi
func (s *Server) adaptMiddleware(h HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// Start démarre le serveur
func (s *Server) Start() error {
	s.setupRoutes()

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.config.Port),
		Handler:           s.router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	log.Printf("🎉 Serveur démarré sur le port %d", s.config.Port)
	log.Printf("📍 URL: http://localhost:%d", s.config.Port)

	if s.config.BasicAuthConfig.Enabled {
		log.Printf("🔒 Authentification Basic Auth activée")
	}

	return s.server.ListenAndServe()
}

// Shutdown arrête proprement le serveur
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Arrêt du serveur...")
	return s.server.Shutdown(ctx)
}

// ensureDir crée un répertoire s'il n'existe pas
func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// fileExists vérifie si un fichier existe
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ensureStaticFiles crée les fichiers statiques par défaut si absents
func EnsureStaticFiles(staticDir string) error {
	// Créer le répertoire des images si nécessaire
	imagesDir := filepath.Join(staticDir, "images")
	if err := ensureDir(imagesDir); err != nil {
		return err
	}

	// Créer le répertoire CSS si nécessaire
	cssDir := filepath.Join(staticDir, "css")
	if err := ensureDir(cssDir); err != nil {
		return err
	}

	// Créer un fichier .gitkeep dans images
	gitkeepPath := filepath.Join(imagesDir, ".gitkeep")
	if !fileExists(gitkeepPath) {
		if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
			return err
		}
	}

	return nil
}
