.PHONY: help run build test lint clean install-tools

help: ## Affiche cette aide
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## Installe les outils de développement
	@echo "Installation des outils..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

deps: ## Télécharge les dépendances
	go mod download
	go mod verify
	go mod tidy

run: ## Lance le serveur en mode dev
	@echo "Démarrage du serveur..."
	@if [ ! -f .env ]; then echo "Copie de .env.example vers .env"; cp configs/.env.example .env; fi
	go run cmd/wedding-web/main.go

watch: ## Lance le serveur avec rechargement automatique (nécessite air)
	@echo "Démarrage en mode watch (hot reload)..."
	@if ! command -v air > /dev/null; then \
		echo "Air n'est pas installé. Installation..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air

build: ## Compile le binaire
	@echo "Compilation..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o wedding-web cmd/wedding-web/main.go

test: ## Lance les tests
	@echo "Exécution des tests..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Lance les tests avec couverture HTML
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Ouvrez coverage.html dans votre navigateur"

lint: ## Lance les linters
	@echo "Linting avec golangci-lint..."
	golangci-lint run --timeout 5m ./...
	@echo "Analyse de sécurité avec gosec..."
	gosec -quiet ./...
	@echo "Analyse statique avec staticcheck..."
	staticcheck ./...

fmt: ## Formate le code
	go fmt ./...
	gofmt -s -w .

vet: ## Vérifie le code avec go vet
	go vet ./...

clean: ## Nettoie les fichiers générés
	rm -f wedding-web
	rm -f coverage.txt coverage.html
	rm -rf rsvp_data/
	go clean -cache

docker-build: ## Construit l'image Docker (optionnel)
	docker build -t wedding-web:latest .

all: deps fmt vet lint test build ## Exécute toutes les vérifications

.DEFAULT_GOAL := help

