# 💍 Site Vitrine de Mariage

Site web élégant et sécurisé pour présenter les informations de votre mariage et gérer les RSVP.

## 🎯 Caractéristiques

- **Architecture hexagonale** propre avec séparation domaine/application/adapters
- **Sécurité renforcée** : chiffrement AES-GCM, CSRF, rate limiting, headers de sécurité stricts
- **Sans base de données** : stockage chiffré dans un fichier JSON
- **Responsive** : interface adaptée à tous les écrans
- **Zero JavaScript** : site fonctionnel sans JS
- **Export calendrier** : téléchargement du planning au format .ics

## 📋 Pages

1. **Page d'accueil** (`/`) - Hero avec photo, présentation et liens principaux
2. **Planning** (`/planning`) - Déroulement de la journée + export .ics
3. **Infos pratiques** (`/infos`) - Lieu, accès, hébergement, dress code
4. **RSVP** (`/rsvp`) - Formulaire de confirmation avec protection anti-spam

## 🏗️ Architecture

```
wedding-web/
├── cmd/wedding-web/        # Point d'entrée de l'application
├── internal/
│   ├── domain/             # Entités métier et ports
│   ├── application/        # Use cases et logique métier
│   └── adapters/           # Implémentations (HTTP, Storage)
├── web/
│   ├── templates/          # Templates HTML
│   └── static/             # CSS et assets
└── configs/                # Fichiers de configuration
```

## 🚀 Installation et démarrage rapide

### Prérequis

- Go 1.22 ou supérieur
- Make (optionnel mais recommandé)

### 1. Cloner le repository

```bash
cd wedding-web
```

### 2. Installer les dépendances

```bash
go mod download
```

Ou avec Make :

```bash
make deps
```

### 3. Configurer les variables d'environnement

Copiez le fichier d'exemple (si absent, les valeurs par défaut seront utilisées) :

```bash
cp .env.example .env
```

Éditez `.env` et **générez une clé de chiffrement** :

```bash
# Générer une clé de 32 bytes en base64
openssl rand -base64 32
```

Copiez la clé générée dans `RSVP_ENCRYPTION_KEY`.

### 4. Démarrer le serveur

```bash
make run
```

Ou directement avec Go :

```bash
go run cmd/wedding-web/main.go
```

Le site sera accessible sur **http://localhost:8080**

## ⚙️ Configuration

### Variables d'environnement

| Variable | Description | Défaut |
|----------|-------------|--------|
| `ENV` | Environnement (dev/prod) | `dev` |
| `PORT` | Port du serveur | `8080` |
| `BASE_URL` | URL de base du site | `http://localhost:8080` |
| `RSVP_ENCRYPTION_KEY` | Clé de chiffrement (32 bytes base64) | ⚠️ **Obligatoire** |
| `RSVP_STORAGE_PATH` | Chemin du fichier de stockage | `./rsvp_data/reservations.json` |
| `BASIC_AUTH_USER` | Utilisateur Basic Auth (optionnel) | - |
| `BASIC_AUTH_PASS` | Mot de passe Basic Auth (optionnel) | - |
| `RATE_LIMIT_PER_MINUTE` | Limite de requêtes par minute | `10` |
| `ENABLE_HSTS` | Activer HSTS (prod uniquement) | `false` |
| `ALLOWED_HOSTS` | Hosts autorisés (séparés par virgule) | `localhost,127.0.0.1` |

### Sécurité : Basic Auth

Pour protéger le site avec Basic Auth (site privé) :

```bash
# Dans .env
BASIC_AUTH_USER=votre_utilisateur
BASIC_AUTH_PASS=votre_mot_de_passe_fort
```

Les visiteurs devront s'authentifier pour accéder au site.

## 🎨 Personnalisation

### 1. Remplacer la photo hero

Placez votre photo dans `web/static/images/hero.jpg`

- Format recommandé : JPG ou PNG
- Dimensions recommandées : 1200x800px minimum
- Poids : < 500KB (optimisez avec tinypng.com ou similaire)

Si aucune photo n'est présente, un placeholder avec emoji sera affiché.

### 2. Modifier les textes

Les textes par défaut sont définis dans :

- **Planning** : `internal/domain/planning.go` → fonction `GetDefaultPlanning()`
- **Infos pratiques** : `internal/domain/info.go` → fonction `GetDefaultPracticalInfo()`

Éditez ces fichiers et recompilez :

```bash
make build
```

### 3. Modifier le style

Le CSS est dans `web/static/css/style.css`. Les variables CSS en début de fichier permettent de changer facilement les couleurs :

```css
:root {
    --primary-color: #d4a574;    /* Couleur principale */
    --primary-dark: #b8935f;      /* Couleur principale foncée */
    --secondary-color: #8b7355;   /* Couleur secondaire */
    /* ... */
}
```

## 🧪 Tests

### Lancer les tests

```bash
make test
```

Ou avec couverture HTML :

```bash
make test-coverage
```

### Linting et analyse statique

```bash
# Installer les outils (une fois)
make install-tools

# Lancer tous les linters
make lint
```

Cela exécute :
- `golangci-lint` - Linting complet
- `gosec` - Analyse de sécurité
- `staticcheck` - Analyse statique

## 🏭 Compilation et déploiement

### Compilation

```bash
make build
```

Cela génère le binaire `wedding-web` à la racine.

### Déploiement en production

1. **Compilez le binaire** sur votre machine ou en CI :

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o wedding-web cmd/wedding-web/main.go
```

2. **Transférez le binaire** et les répertoires `web/` sur votre serveur

3. **Configurez les variables d'environnement** :

```bash
export ENV=prod
export PORT=8080
export RSVP_ENCRYPTION_KEY="votre-clé-générée-32-bytes-base64"
export ENABLE_HSTS=true
export ALLOWED_HOSTS="votredomaine.com,www.votredomaine.com"
```

4. **Lancez le serveur** :

```bash
./wedding-web
```

### Reverse proxy (Nginx ou Traefik)

Le serveur Go écoute sur le port configuré (8080 par défaut). Configurez votre reverse proxy pour router le trafic HTTPS vers ce port.

**Exemple Nginx** :

```nginx
server {
    listen 443 ssl http2;
    server_name votredomaine.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Exemple Traefik (docker-compose)** :

```yaml
services:
  wedding-web:
    image: wedding-web:latest
    environment:
      - ENV=prod
      - RSVP_ENCRYPTION_KEY=${RSVP_ENCRYPTION_KEY}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.wedding.rule=Host(`votredomaine.com`)"
      - "traefik.http.routers.wedding.tls=true"
      - "traefik.http.routers.wedding.tls.certresolver=letsencrypt"
```

## 🔒 Sécurité

Le projet implémente de nombreuses mesures de sécurité :

### Chiffrement
- ✅ Stockage RSVP chiffré avec AES-256-GCM
- ✅ Clé de chiffrement en variable d'environnement

### Protection des formulaires
- ✅ Tokens CSRF sur tous les POST
- ✅ Honeypot anti-spam
- ✅ Rate limiting par IP
- ✅ Validation serveur stricte
- ✅ Taille maximale du body (1MB)

### Headers HTTP
- ✅ Content-Security-Policy stricte
- ✅ X-Content-Type-Options: nosniff
- ✅ X-Frame-Options: DENY
- ✅ Referrer-Policy
- ✅ Permissions-Policy
- ✅ HSTS (en production)

### Autres
- ✅ Basic Auth optionnelle (comparaison en temps constant)
- ✅ Pas de logs PII (informations personnelles)
- ✅ Gestion propre des erreurs (pas de stacktrace en prod)
- ✅ Timeouts HTTP configurés
- ✅ Shutdown graceful

### Audit de sécurité

```bash
make lint
```

Cela lance `gosec` qui analyse le code pour détecter les vulnérabilités.

## 📝 Structure des données RSVP

Les RSVP sont stockés dans un fichier JSON chiffré :

```json
{
  "rsvps": [
    {
      "id": "...",
      "first_name": "Jean",
      "last_name": "Dupont",
      "adults_count": 2,
      "children_count": 1,
      "allergies": "Végétarien",
      "message": "Hâte d'être là !",
      "submitted_at": "2025-03-15T10:30:00Z"
    }
  ]
}
```

**Note** : L'IP du visiteur n'est pas persistée pour respecter la vie privée.

## 🛠️ Commandes Make disponibles

```bash
make help              # Affiche l'aide
make deps              # Télécharge les dépendances
make run               # Lance le serveur en dev
make build             # Compile le binaire
make test              # Lance les tests
make test-coverage     # Tests avec couverture HTML
make lint              # Lance tous les linters
make fmt               # Formate le code
make vet               # go vet
make clean             # Nettoie les fichiers générés
make install-tools     # Installe les outils de dev
```

## 📦 Dépendances

Le projet utilise un minimum de dépendances :

- **github.com/go-chi/chi/v5** - Router HTTP léger et compatible stdlib
- **golang.org/x/time** - Rate limiting

Aucune dépendance lourde, tout est conçu pour être simple et maintenable.

## 🤝 Contribution

Ce projet est prévu pour un usage personnel (site de mariage), mais si vous souhaitez l'améliorer :

1. Fork le projet
2. Créez une branche (`git checkout -b feature/amelioration`)
3. Committez vos changements (`git commit -am 'Ajout fonctionnalité'`)
4. Pushez (`git push origin feature/amelioration`)
5. Ouvrez une Pull Request

## 📄 Licence

Ce projet est fourni "tel quel" sans garantie. Libre à vous de l'utiliser et le modifier pour votre mariage ! 💕

## 💡 Inspirations

Architecture inspirée du repository [github.com/AF1890/profil](https://github.com/AF1890/profil) avec une séparation claire domaine/application/adapters.

## 🎉 Bon mariage !

Profitez de cette belle journée et que ce site vous aide à l'organiser sereinement ! 💍

---

**Fait avec ❤️ en Go**

