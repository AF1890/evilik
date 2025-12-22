package http

import (
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// HandlerFunc définit le type de handler
type HandlerFunc func(ResponseWriter, *Request) error

// ResponseWriter wrapper pour capturer le status code
type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Written() int
}

type responseWriter struct {
	http.ResponseWriter
	status  int
	written int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

func (w *responseWriter) Status() int {
	if w.status == 0 {
		return 200
	}
	return w.status
}

func (w *responseWriter) Written() int {
	return w.written
}

// Request wrapper
type Request struct {
	*http.Request
	requestID string
}

func (r *Request) RequestID() string {
	return r.requestID
}

// Middleware type
type Middleware func(HandlerFunc) HandlerFunc

// requestIDMiddleware ajoute un ID unique à chaque requête
func requestIDMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			// Générer un ID unique
			b := make([]byte, 16)
			rand.Read(b)
			requestID := hex.EncodeToString(b)

			r.requestID = requestID
			w.Header().Set("X-Request-ID", requestID)

			return next(w, r)
		}
	}
}

// loggingMiddleware log les requêtes
func loggingMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			start := time.Now()

			err := next(w, r)

			duration := time.Since(start)
			status := w.Status()

			// Ne pas logger les informations sensibles (PII)
			log.Printf("[%s] %s %s %d %dms (request_id=%s)",
				r.Method,
				r.URL.Path,
				getClientIP(r.Request),
				status,
				duration.Milliseconds(),
				r.RequestID(),
			)

			return err
		}
	}
}

// recoverMiddleware récupère les panics
func recoverMiddleware(isProd bool) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) (err error) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.Printf("PANIC: %v (request_id=%s)\n%s", rvr, r.RequestID(), string(debug.Stack()))

					w.WriteHeader(http.StatusInternalServerError)

					if isProd {
						w.Write([]byte("Une erreur interne s'est produite"))
					} else {
						w.Write([]byte(fmt.Sprintf("Panic: %v\n\n%s", rvr, string(debug.Stack()))))
					}
				}
			}()

			return next(w, r)
		}
	}
}

// securityHeadersMiddleware ajoute les headers de sécurité
func securityHeadersMiddleware(enableHSTS bool) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			// CSP - autoriser le script inline pour le formulaire RSVP uniquement
			w.Header().Set("Content-Security-Policy",
				"default-src 'self'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data: https:; "+
					"font-src 'self'; "+
					"script-src 'self' 'unsafe-inline'; "+
					"frame-ancestors 'none'; "+
					"base-uri 'self'; "+
					"form-action 'self'")

			// Autres headers de sécurité
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

			// HSTS en production uniquement
			if enableHSTS {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}

			return next(w, r)
		}
	}
}

// gzipMiddleware compresse les réponses
func gzipMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			// Vérifier si le client supporte gzip
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				return next(w, r)
			}

			// Créer le writer gzip
			gz := gzip.NewWriter(w)
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")

			gzw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}

			return next(gzw, r)
		}
	}
}

type gzipResponseWriter struct {
	ResponseWriter
	io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// cacheControlMiddleware définit les politiques de cache
func cacheControlMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			// Cache pour les assets statiques
			if strings.HasPrefix(r.URL.Path, "/static/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				// Pas de cache pour les pages dynamiques
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			}

			return next(w, r)
		}
	}
}

// maxBytesMiddleware limite la taille du body
func maxBytesMiddleware(maxBytes int64) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			return next(w, r)
		}
	}
}

// timeoutMiddleware ajoute un timeout aux requêtes
func timeoutMiddleware(timeout time.Duration) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r.Request = r.Request.WithContext(ctx)

			return next(w, r)
		}
	}
}

// getClientIP récupère l'IP réelle du client
func getClientIP(r *http.Request) string {
	// Vérifier X-Forwarded-For (proxy/load balancer)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Prendre la première IP
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Vérifier X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Sinon, utiliser RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// Chain chaîne plusieurs middlewares
func Chain(handler HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
