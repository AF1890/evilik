#!/bin/bash
set -e

echo "🚀 Installation du site de mariage sur VPS"
echo "==========================================="
echo ""

# Couleurs pour les messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Variables à configurer
DOMAIN="votre-domaine.fr"  # À MODIFIER
APP_USER="wedding"
APP_DIR="/opt/wedding-web"
GO_VERSION="1.22.0"

echo -e "${BLUE}📋 Configuration :${NC}"
echo "  - Domaine : $DOMAIN"
echo "  - Utilisateur : $APP_USER"
echo "  - Répertoire : $APP_DIR"
echo ""

# Vérifier si on est root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}❌ Ce script doit être exécuté en tant que root${NC}"
    exit 1
fi

# 1. Mise à jour du système
echo -e "${GREEN}📦 1. Mise à jour du système...${NC}"
apt update && apt upgrade -y

# 2. Installation des dépendances
echo -e "${GREEN}📦 2. Installation des dépendances...${NC}"
apt install -y nginx certbot python3-certbot-nginx git curl wget ufw

# 3. Installation de Go
echo -e "${GREEN}🐹 3. Installation de Go $GO_VERSION...${NC}"
if ! command -v go &> /dev/null; then
    cd /tmp
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
    rm go${GO_VERSION}.linux-amd64.tar.gz
    echo -e "${GREEN}✅ Go installé : $(go version)${NC}"
else
    echo -e "${GREEN}✅ Go déjà installé : $(go version)${NC}"
fi

# 4. Création de l'utilisateur pour l'application
echo -e "${GREEN}👤 4. Création de l'utilisateur $APP_USER...${NC}"
if ! id "$APP_USER" &>/dev/null; then
    useradd -r -s /bin/bash -d $APP_DIR $APP_USER
    echo -e "${GREEN}✅ Utilisateur $APP_USER créé${NC}"
else
    echo -e "${GREEN}✅ Utilisateur $APP_USER existe déjà${NC}"
fi

# 5. Création des répertoires
echo -e "${GREEN}📁 5. Création des répertoires...${NC}"
mkdir -p $APP_DIR/{rsvp_data,logs,backup}
chown -R $APP_USER:$APP_USER $APP_DIR

# 6. Configuration du firewall
echo -e "${GREEN}🔥 6. Configuration du firewall (UFW)...${NC}"
ufw --force enable
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 'Nginx Full'
ufw --force reload
echo -e "${GREEN}✅ Firewall configuré${NC}"

# 7. Configuration Nginx
echo -e "${GREEN}🌐 7. Configuration de Nginx...${NC}"
cat > /etc/nginx/sites-available/$DOMAIN <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name $DOMAIN www.$DOMAIN;

    # Redirect to HTTPS (will be uncommented after certbot)
    # return 301 https://\$server_name\$request_uri;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }

    # Serve static files directly
    location /static/ {
        alias $APP_DIR/web/static/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
EOF

ln -sf /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t && systemctl reload nginx
echo -e "${GREEN}✅ Nginx configuré${NC}"

# 8. Création du fichier .env de production
echo -e "${GREEN}⚙️  8. Création du fichier de configuration...${NC}"
cat > $APP_DIR/.env <<EOF
# Configuration de production
ENV=prod
PORT=8080
BASE_URL=https://$DOMAIN

# Génération d'une clé de chiffrement forte
RSVP_ENCRYPTION_KEY=$(openssl rand -base64 32)

# Admin credentials (À MODIFIER !)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=$(openssl rand -base64 16)

# Paths
RSVP_STORAGE_PATH=$APP_DIR/rsvp_data/reservations.json

# Config
CONFIG_PATH=$APP_DIR/conf/prod.yaml
EOF

chown $APP_USER:$APP_USER $APP_DIR/.env
chmod 600 $APP_DIR/.env

echo -e "${BLUE}📝 Credentials admin générés (à noter) :${NC}"
echo "  Username: admin"
echo "  Password: $(grep ADMIN_PASSWORD $APP_DIR/.env | cut -d= -f2)"
echo ""

# 9. Création du service systemd
echo -e "${GREEN}🔧 9. Configuration du service systemd...${NC}"
cat > /etc/systemd/system/wedding-web.service <<EOF
[Unit]
Description=Wedding Web Application
After=network.target

[Service]
Type=simple
User=$APP_USER
WorkingDirectory=$APP_DIR
EnvironmentFile=$APP_DIR/.env
ExecStart=$APP_DIR/wedding-web
Restart=always
RestartSec=10

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$APP_DIR/rsvp_data $APP_DIR/logs

# Logging
StandardOutput=append:$APP_DIR/logs/app.log
StandardError=append:$APP_DIR/logs/error.log

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable wedding-web.service
echo -e "${GREEN}✅ Service systemd configuré${NC}"

# 10. Script de backup automatique
echo -e "${GREEN}💾 10. Configuration du backup automatique...${NC}"
cat > $APP_DIR/backup.sh <<'BACKUP_SCRIPT'
#!/bin/bash
BACKUP_DIR="/opt/wedding-web/backup"
RSVP_FILE="/opt/wedding-web/rsvp_data/reservations.json"
DATE=$(date +%Y%m%d_%H%M%S)

if [ -f "$RSVP_FILE" ]; then
    cp "$RSVP_FILE" "$BACKUP_DIR/reservations_$DATE.json"
    # Garder seulement les 30 derniers backups
    ls -t $BACKUP_DIR/reservations_*.json | tail -n +31 | xargs -r rm
    echo "$(date): Backup créé - reservations_$DATE.json" >> $BACKUP_DIR/backup.log
fi
BACKUP_SCRIPT

chmod +x $APP_DIR/backup.sh
chown $APP_USER:$APP_USER $APP_DIR/backup.sh

# Ajouter au cron (tous les jours à 2h du matin)
(crontab -u $APP_USER -l 2>/dev/null; echo "0 2 * * * $APP_DIR/backup.sh") | crontab -u $APP_USER -
echo -e "${GREEN}✅ Backup automatique configuré (tous les jours à 2h)${NC}"

echo ""
echo -e "${GREEN}=============================================${NC}"
echo -e "${GREEN}✅ Installation de base terminée !${NC}"
echo -e "${GREEN}=============================================${NC}"
echo ""
echo -e "${BLUE}📋 Prochaines étapes :${NC}"
echo ""
echo "1. Déployez votre application :"
echo "   ${BLUE}scp -r /chemin/local/wedding-web root@IP_VPS:$APP_DIR/${NC}"
echo ""
echo "2. Compilez l'application sur le VPS :"
echo "   ${BLUE}cd $APP_DIR && go build -o wedding-web cmd/wedding-web/*.go${NC}"
echo ""
echo "3. Démarrez le service :"
echo "   ${BLUE}systemctl start wedding-web${NC}"
echo ""
echo "4. Configurez HTTPS avec Let's Encrypt :"
echo "   ${BLUE}certbot --nginx -d $DOMAIN -d www.$DOMAIN${NC}"
echo ""
echo "5. Vérifiez le statut :"
echo "   ${BLUE}systemctl status wedding-web${NC}"
echo "   ${BLUE}tail -f $APP_DIR/logs/app.log${NC}"
echo ""
echo -e "${RED}⚠️  N'oubliez pas de :${NC}"
echo "  - Noter les credentials admin ci-dessus"
echo "  - Pointer votre domaine vers l'IP du VPS"
echo "  - Modifier le domaine dans ce script"
echo ""

