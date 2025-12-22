package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config contient toute la configuration de l'application.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Security SecurityConfig `yaml:"security"`
	RSVP     RSVPConfig     `yaml:"rsvp"`
	Admin    AdminConfig    `yaml:"admin"`
}

// ServerConfig contient la configuration du serveur HTTP.
type ServerConfig struct {
	Port         string `yaml:"port"`
	BaseURL      string `yaml:"base_url"`
	Environment  string `yaml:"environment"` // dev, prod
	AllowedHosts string `yaml:"allowed_hosts"`
}

// SecurityConfig contient la configuration de sécurité.
type SecurityConfig struct {
	CSPEnabled          bool   `yaml:"csp_enabled"`
	HSTSEnabled         bool   `yaml:"hsts_enabled"`
	RateLimitEnabled    bool   `yaml:"rate_limit_enabled"`
	RateLimitPerMinute  int    `yaml:"rate_limit_per_minute"`
	BasicAuthEnabled    bool   `yaml:"basic_auth_enabled"`
	SessionSecret       string `yaml:"session_secret"`
	EncryptionKey       string `yaml:"encryption_key"`
	EncryptionKeyEnvVar string `yaml:"encryption_key_env_var"`
}

// RSVPConfig contient la configuration du système RSVP.
type RSVPConfig struct {
	Enabled     bool   `yaml:"enabled"`
	StoragePath string `yaml:"storage_path"`
}

// AdminConfig contient la configuration de la page admin.
type AdminConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	UsernameEnvVar string `yaml:"username_env_var"`
	PasswordEnvVar string `yaml:"password_env_var"`
}

// Defaults définit les valeurs par défaut de la configuration.
func (c *Config) Defaults() {
	// Server defaults
	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Server.Environment == "" {
		c.Server.Environment = "dev"
	}
	if c.Server.BaseURL == "" {
		c.Server.BaseURL = fmt.Sprintf("http://localhost:%s", c.Server.Port)
	}

	// Security defaults
	if c.Security.RateLimitPerMinute == 0 {
		c.Security.RateLimitPerMinute = 60
	}
	if c.Security.SessionSecret == "" {
		c.Security.SessionSecret = "default-dev-session-secret-change-in-prod"
	}

	// RSVP defaults
	if c.RSVP.StoragePath == "" {
		c.RSVP.StoragePath = "./rsvp_data/reservations.json"
	}
}

// LoadFromEnv charge les secrets depuis les variables d'environnement.
func (c *Config) LoadFromEnv() error {
	// Charger la clé de chiffrement depuis ENV si spécifiée
	if c.Security.EncryptionKeyEnvVar != "" {
		if key := os.Getenv(c.Security.EncryptionKeyEnvVar); key != "" {
			c.Security.EncryptionKey = key
		}
	}

	// Charger les credentials admin depuis ENV si spécifiés
	if c.Admin.UsernameEnvVar != "" {
		if username := os.Getenv(c.Admin.UsernameEnvVar); username != "" {
			c.Admin.Username = username
		}
	}
	if c.Admin.PasswordEnvVar != "" {
		if password := os.Getenv(c.Admin.PasswordEnvVar); password != "" {
			c.Admin.Password = password
		}
	}

	// Override avec PORT si défini
	if port := os.Getenv("PORT"); port != "" {
		c.Server.Port = port
	}

	// Override avec ENV si défini
	if env := os.Getenv("ENV"); env != "" {
		c.Server.Environment = env
	}

	return nil
}

// Validate valide la configuration.
func (c *Config) Validate() error {
	// En production, la clé de chiffrement est obligatoire
	if c.IsProd() {
		if c.Security.EncryptionKey == "" {
			return fmt.Errorf("encryption_key est obligatoire en production")
		}
		if len(c.Security.EncryptionKey) < 32 {
			return fmt.Errorf("encryption_key doit faire au moins 32 caractères en production")
		}
	}

	// Si admin est activé, username et password sont obligatoires
	if c.Admin.Enabled {
		if c.Admin.Username == "" || c.Admin.Password == "" {
			return fmt.Errorf("admin.username et admin.password sont obligatoires si admin est activé")
		}
		if c.IsProd() && c.Admin.Password == "changeme" {
			return fmt.Errorf("vous devez changer le mot de passe admin en production")
		}
	}

	// Générer une clé temporaire en dev si nécessaire
	if c.IsDev() && c.Security.EncryptionKey == "" {
		log.Println("⚠️  ATTENTION: Aucune clé de chiffrement configurée !")
		log.Println("   Générez une clé avec: openssl rand -base64 32")
		log.Println("   Puis ajoutez-la dans RSVP_ENCRYPTION_KEY")
		log.Println("   Pour le développement, une clé temporaire sera générée.")
		c.Security.EncryptionKey = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
		log.Println("⚠️  Clé temporaire générée (DÉVELOPPEMENT UNIQUEMENT)")
	}

	// Afficher le statut de l'admin
	if c.Admin.Enabled && c.Admin.Username != "" && c.Admin.Password != "" {
		log.Printf("✅ Page admin activée sur /admin (utilisateur: %s)", c.Admin.Username)
	} else {
		log.Println("⚠️  Page admin désactivée")
	}

	return nil
}

// IsDev retourne true si l'environnement est dev.
func (c *Config) IsDev() bool {
	return c.Server.Environment == "dev"
}

// IsProd retourne true si l'environnement est prod.
func (c *Config) IsProd() bool {
	return c.Server.Environment == "prod"
}

// LoadConfig charge la configuration depuis un fichier YAML.
func LoadConfig(configPath string) (*Config, error) {
	// Lire le fichier YAML
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture fichier config %s: %w", configPath, err)
	}

	// Parser le YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("erreur parsing YAML: %w", err)
	}

	// Appliquer les valeurs par défaut
	cfg.Defaults()

	// Charger les secrets depuis les variables d'environnement
	if err := cfg.LoadFromEnv(); err != nil {
		return nil, fmt.Errorf("erreur chargement variables d'environnement: %w", err)
	}

	// Valider la configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("erreur validation config: %w", err)
	}

	return &cfg, nil
}

// GetConfigPath retourne le chemin du fichier de configuration basé sur ENV.
func GetConfigPath() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		return configPath
	}

	return fmt.Sprintf("conf/%s.yaml", env)
}
