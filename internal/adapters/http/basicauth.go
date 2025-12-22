package http

import (
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

// BasicAuthConfig contient la configuration pour l'authentification basique
type BasicAuthConfig struct {
	Username string
	Password string
	Enabled  bool
}

// basicAuthMiddleware vérifie l'authentification Basic Auth
func basicAuthMiddleware(config BasicAuthConfig) func(next HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) error {
			// Si Basic Auth n'est pas activé, continuer
			if !config.Enabled {
				return next(w, r)
			}

			// Récupérer l'en-tête Authorization
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Wedding Site"`)
				w.WriteHeader(401)
				w.Write([]byte("Authentification requise"))
				return nil
			}

			// Vérifier le format "Basic <credentials>"
			const prefix = "Basic "
			if !strings.HasPrefix(auth, prefix) {
				w.WriteHeader(401)
				w.Write([]byte("Format d'authentification invalide"))
				return nil
			}

			// Décoder les credentials
			decoded, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
			if err != nil {
				w.WriteHeader(401)
				w.Write([]byte("Credentials invalides"))
				return nil
			}

			// Séparer username:password
			credentials := string(decoded)
			colonIndex := strings.Index(credentials, ":")
			if colonIndex == -1 {
				w.WriteHeader(401)
				w.Write([]byte("Format credentials invalide"))
				return nil
			}

			username := credentials[:colonIndex]
			password := credentials[colonIndex+1:]

			// Comparaison en temps constant pour éviter les timing attacks
			usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(config.Username)) == 1
			passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(config.Password)) == 1

			if !usernameMatch || !passwordMatch {
				w.Header().Set("WWW-Authenticate", `Basic realm="Wedding Site"`)
				w.WriteHeader(401)
				w.Write([]byte("Authentification échouée"))
				return nil
			}

			return next(w, r)
		}
	}
}
