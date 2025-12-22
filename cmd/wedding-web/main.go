package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"wedding-web/internal/adapters/http"
	"wedding-web/internal/adapters/storage"
	"wedding-web/internal/application"
)

func main() {
	log.Println("🎉 Démarrage de l'application Wedding Web...")

	// Charger la configuration depuis le fichier YAML
	configPath := GetConfigPath()
	log.Printf("📄 Chargement de la configuration depuis: %s", configPath)

	appConfig, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la configuration: %v", err)
	}

	log.Printf("🌍 Environnement: %s", appConfig.Server.Environment)
	log.Printf("🚪 Port: %s", appConfig.Server.Port)

	// Initialiser les services
	services, err := initializeServices(appConfig)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation des services: %v", err)
	}

	// Déterminer les répertoires
	staticDir := getEnv("STATIC_DIR", "./web/static")
	templatesDir := getEnv("TEMPLATES_DIR", "./web/templates")

	// Créer les handlers
	handlers, err := http.NewHandlers(
		services.rsvpService,
		services.planningService,
		services.infoService,
		services.calendarService,
		services.csrfManager,
		templatesDir,
		appConfig.IsDev(),
		appConfig.Admin.Username,
		appConfig.Admin.Password,
	)
	if err != nil {
		log.Fatalf("Erreur lors de la création des handlers: %v", err)
	}

	// Créer le serveur
	port, err := strconv.Atoi(appConfig.Server.Port)
	if err != nil {
		log.Fatalf("Port invalide: %v", err)
	}

	serverConfig := http.ServerConfig{
		Port:       port,
		IsProd:     appConfig.IsProd(),
		EnableHSTS: appConfig.Security.HSTSEnabled,
		BasicAuthConfig: http.BasicAuthConfig{
			Username: "", // Basic auth global désactivé par défaut
			Password: "",
			Enabled:  false,
		},
		RateLimitPerMinute: appConfig.Security.RateLimitPerMinute,
		MaxBodySize:        1 << 20, // 1 MB
		StaticDir:          staticDir,
		TemplatesDir:       templatesDir,
	}

	server := http.NewServer(serverConfig, handlers)

	// Assurer que les répertoires statiques existent
	if err := http.EnsureStaticFiles(staticDir); err != nil {
		log.Fatalf("Erreur lors de la création des répertoires statiques: %v", err)
	}

	// Canal pour les signaux système
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Démarrer le serveur dans une goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()

	// Attendre le signal d'arrêt
	<-sigChan

	// Arrêt propre du serveur
	log.Println("Signal d'arrêt reçu, fermeture du serveur...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Erreur lors de l'arrêt du serveur: %v", err)
	}

	log.Println("Serveur arrêté proprement. Au revoir ! 👋")
}

// Services contient tous les services de l'application
type Services struct {
	rsvpService     *application.RSVPService
	planningService *application.PlanningService
	infoService     *application.InfoService
	calendarService *application.CalendarService
	csrfManager     *http.CSRFManager
}

// initializeServices initialise tous les services
func initializeServices(config *Config) (*Services, error) {
	// Storage pour les RSVP
	rsvpStorage, err := storage.NewEncryptedFileStorage(
		config.RSVP.StoragePath,
		config.Security.EncryptionKey,
	)
	if err != nil {
		return nil, err
	}

	// Services métier
	rsvpService := application.NewRSVPService(rsvpStorage)
	planningService := application.NewPlanningService()
	infoService := application.NewInfoService()
	calendarService := application.NewCalendarService()

	// CSRF Manager
	csrfManager := http.NewCSRFManager()

	return &Services{
		rsvpService:     rsvpService,
		planningService: planningService,
		infoService:     infoService,
		calendarService: calendarService,
		csrfManager:     csrfManager,
	}, nil
}

// getEnv récupère une variable d'environnement avec une valeur par défaut
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
