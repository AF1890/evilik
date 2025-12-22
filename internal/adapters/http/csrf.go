package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"net/http"
	"sync"
	"time"
)

// CSRFToken représente un token CSRF
type CSRFToken struct {
	Value     string
	ExpiresAt time.Time
}

// CSRFManager gère les tokens CSRF
type CSRFManager struct {
	tokens map[string]CSRFToken
	mu     sync.RWMutex
}

// NewCSRFManager crée un nouveau gestionnaire CSRF
func NewCSRFManager() *CSRFManager {
	cm := &CSRFManager{
		tokens: make(map[string]CSRFToken),
	}

	// Nettoyer les tokens expirés toutes les 10 minutes
	go cm.cleanupExpiredTokens()

	return cm
}

// GenerateToken génère un nouveau token CSRF pour une session
func (cm *CSRFManager) GenerateToken(sessionID string) (string, error) {
	// Générer un token aléatoire
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	tokenValue := base64.URLEncoding.EncodeToString(b)

	token := CSRFToken{
		Value:     tokenValue,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	cm.mu.Lock()
	cm.tokens[sessionID] = token
	cm.mu.Unlock()

	return tokenValue, nil
}

// ValidateToken valide un token CSRF
func (cm *CSRFManager) ValidateToken(sessionID, tokenValue string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	token, exists := cm.tokens[sessionID]
	if !exists {
		return false
	}

	// Vérifier l'expiration
	if time.Now().After(token.ExpiresAt) {
		return false
	}

	// Comparer le token
	return token.Value == tokenValue
}

// cleanupExpiredTokens nettoie les tokens expirés
func (cm *CSRFManager) cleanupExpiredTokens() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cm.mu.Lock()
		now := time.Now()
		for sessionID, token := range cm.tokens {
			if now.After(token.ExpiresAt) {
				delete(cm.tokens, sessionID)
			}
		}
		cm.mu.Unlock()
	}
}

// getOrCreateSession récupère ou crée une session
func getOrCreateSession(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Créer une nouvelle session
	b := make([]byte, 32)
	rand.Read(b)
	sessionID := base64.URLEncoding.EncodeToString(b)

	// Créer le cookie de session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   3600, // 1 heure
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return sessionID
}

func init() {
	// Nécessaire pour l'encodage/décodage des sessions
	gob.Register(CSRFToken{})
}
