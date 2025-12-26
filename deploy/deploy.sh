#!/bin/bash
set -e

echo "📦 Déploiement de l'application Wedding Web"
echo "============================================"
echo ""

# Variables
APP_DIR="/opt/wedding-web"
VPS_IP="VOTRE_IP_VPS"  # À MODIFIER
DOMAIN="votre-domaine.fr"  # À MODIFIER

# Couleurs
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}📋 Configuration :${NC}"
echo "  - VPS IP : $VPS_IP"
echo "  - Domaine : $DOMAIN"
echo ""

# Vérification des variables
if [ "$VPS_IP" = "VOTRE_IP_VPS" ]; then
    echo -e "${RED}❌ Veuillez modifier VPS_IP dans le script${NC}"
    exit 1
fi

# 1. Compilation locale (optionnel, on peut compiler sur le VPS)
echo -e "${GREEN}🔨 1. Compilation de l'application...${NC}"
GOOS=linux GOARCH=amd64 go build -o wedding-web cmd/wedding-web/*.go
echo -e "${GREEN}✅ Application compilée${NC}"

# 2. Création de l'archive
echo -e "${GREEN}📦 2. Création de l'archive...${NC}"
tar -czf wedding-web.tar.gz \
    wedding-web \
    web/ \
    conf/ \
    go.mod \
    go.sum \
    README.md

echo -e "${GREEN}✅ Archive créée : wedding-web.tar.gz${NC}"

# 3. Upload vers le VPS
echo -e "${GREEN}📤 3. Upload vers le VPS...${NC}"
scp wedding-web.tar.gz root@$VPS_IP:/tmp/
echo -e "${GREEN}✅ Upload terminé${NC}"

# 4. Déploiement sur le VPS
echo -e "${GREEN}🚀 4. Déploiement sur le VPS...${NC}"
ssh root@$VPS_IP << 'ENDSSH'
    cd /tmp
    tar -xzf wedding-web.tar.gz -C /opt/wedding-web/
    chown -R wedding:wedding /opt/wedding-web
    
    # Redémarrage du service
    systemctl restart wedding-web
    systemctl status wedding-web --no-pager
    
    # Nettoyage
    rm /tmp/wedding-web.tar.gz
ENDSSH

echo -e "${GREEN}✅ Déploiement terminé !${NC}"

# 5. Vérification
echo -e "${GREEN}🔍 5. Vérification...${NC}"
sleep 3
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://$VPS_IP:8080/

echo ""
echo -e "${GREEN}=============================================${NC}"
echo -e "${GREEN}✅ Déploiement réussi !${NC}"
echo -e "${GREEN}=============================================${NC}"
echo ""
echo -e "${BLUE}🌐 Accès :${NC}"
echo "  - HTTP : http://$DOMAIN"
echo "  - Admin : http://$DOMAIN/admin"
echo ""
echo -e "${BLUE}📋 Commandes utiles :${NC}"
echo "  - Logs : ${BLUE}ssh root@$VPS_IP 'tail -f /opt/wedding-web/logs/app.log'${NC}"
echo "  - Status : ${BLUE}ssh root@$VPS_IP 'systemctl status wedding-web'${NC}"
echo "  - Restart : ${BLUE}ssh root@$VPS_IP 'systemctl restart wedding-web'${NC}"
echo ""

# Nettoyage local
rm wedding-web wedding-web.tar.gz

