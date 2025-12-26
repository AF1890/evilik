# 🚀 Guide de déploiement - Site de mariage

## 📋 Prérequis

- ✅ VPS Infomaniak créé (Ubuntu 24.04)
- ✅ Clé SSH configurée
- ✅ Nom de domaine pointant vers l'IP du VPS

---

## 🎯 Étape 1 : Configuration DNS

### Chez votre registrar (OVH, Gandi, etc.)

Ajoutez ces enregistrements DNS :

```
Type    Nom     Valeur
A       @       VOTRE_IP_VPS
A       www     VOTRE_IP_VPS
```

**Temps de propagation** : 5 minutes à 24h

---

## 🔧 Étape 2 : Installation sur le VPS

### 2.1 Connexion au VPS

```bash
ssh root@VOTRE_IP_VPS
```

### 2.2 Téléchargement du script d'installation

```bash
# Sur votre machine locale
cd /home/aylin.devillechenous/Documents/Documents\ -\ personnel\ -\ prive/evlilik
scp deploy/setup.sh root@VOTRE_IP_VPS:/root/
```

### 2.3 Modification du script

Sur le VPS :

```bash
nano /root/setup.sh
```

**Modifiez la ligne 16 :**
```bash
DOMAIN="votre-domaine.fr"  # ⬅️ Remplacez par votre domaine
```

### 2.4 Exécution du script d'installation

```bash
chmod +x /root/setup.sh
./setup.sh
```

**Ce script va :**
- ✅ Installer Go, Nginx, Certbot
- ✅ Créer l'utilisateur `wedding`
- ✅ Configurer le firewall (UFW)
- ✅ Générer les clés de chiffrement
- ✅ Créer le service systemd
- ✅ Configurer le backup automatique

⏱️ **Durée** : ~5-10 minutes

---

## 📦 Étape 3 : Déploiement de l'application

### 3.1 Sur votre machine locale

```bash
cd /home/aylin.devillechenous/Documents/Documents\ -\ personnel\ -\ prive/evlilik
nano deploy/deploy.sh
```

**Modifiez les lignes 9-10 :**
```bash
VPS_IP="123.45.67.89"     # ⬅️ IP de votre VPS
DOMAIN="votre-domaine.fr"  # ⬅️ Votre domaine
```

### 3.2 Lancement du déploiement

```bash
chmod +x deploy/deploy.sh
./deploy/deploy.sh
```

**Ce script va :**
- ✅ Compiler l'application pour Linux
- ✅ Créer une archive
- ✅ L'envoyer sur le VPS
- ✅ Déployer et redémarrer le service

---

## 🔒 Étape 4 : Configuration HTTPS (Let's Encrypt)

### 4.1 Sur le VPS

```bash
ssh root@VOTRE_IP_VPS
certbot --nginx -d votre-domaine.fr -d www.votre-domaine.fr
```

**Questions de Certbot :**
- Email : Entrez votre email
- Termes : Acceptez
- Partager l'email : Non (optionnel)
- Redirect HTTP → HTTPS : **Oui** (recommandé)

✅ Certificat valide 90 jours (renouvellement automatique)

---

## ✅ Étape 5 : Vérification

### 5.1 Vérifier que le service tourne

```bash
ssh root@VOTRE_IP_VPS
systemctl status wedding-web
```

**Devrait afficher** : `● wedding-web.service - Wedding Web Application`  
**État** : `active (running)`

### 5.2 Vérifier les logs

```bash
tail -f /opt/wedding-web/logs/app.log
```

### 5.3 Tester le site

Ouvrez votre navigateur :
- **Site** : `https://votre-domaine.fr`
- **Admin** : `https://votre-domaine.fr/admin`

---

## 🔑 Credentials Admin

Les credentials sont générés automatiquement lors de l'installation.

**Pour les récupérer :**

```bash
ssh root@VOTRE_IP_VPS
cat /opt/wedding-web/.env | grep ADMIN
```

**Vous verrez :**
```
ADMIN_USERNAME=admin
ADMIN_PASSWORD=VotreMo...
```

⚠️ **Notez-les précieusement !**

---

## 🔄 Mise à jour de l'application

### Après avoir modifié le code :

```bash
cd /home/aylin.devillechenous/Documents/Documents\ -\ personnel\ -\ prive/evlilik
./deploy/deploy.sh
```

C'est tout ! Le script redéploie automatiquement.

---

## 💾 Backup des données

### Backup automatique

✅ **Configuré automatiquement** : tous les jours à 2h du matin

**Localisation** : `/opt/wedding-web/backup/`

### Backup manuel

```bash
# Télécharger les RSVPs
scp root@VOTRE_IP_VPS:/opt/wedding-web/rsvp_data/reservations.json ./backup/

# Ou via l'interface admin
# https://votre-domaine.fr/admin → Exporter en Excel
```

---

## 🛠️ Commandes utiles

### Sur le VPS

```bash
# Voir les logs en temps réel
tail -f /opt/wedding-web/logs/app.log

# Redémarrer l'application
systemctl restart wedding-web

# Voir le statut
systemctl status wedding-web

# Voir les backups
ls -lh /opt/wedding-web/backup/

# Vérifier Nginx
nginx -t
systemctl status nginx
```

### Depuis votre machine locale

```bash
# Se connecter au VPS
ssh root@VOTRE_IP_VPS

# Voir les logs à distance
ssh root@VOTRE_IP_VPS 'tail -f /opt/wedding-web/logs/app.log'

# Redémarrer à distance
ssh root@VOTRE_IP_VPS 'systemctl restart wedding-web'
```

---

## 🔍 Dépannage

### Le site ne s'affiche pas

```bash
# Vérifier que le service tourne
systemctl status wedding-web

# Vérifier les logs
tail -50 /opt/wedding-web/logs/error.log

# Vérifier Nginx
nginx -t
systemctl status nginx

# Vérifier le firewall
ufw status
```

### Erreur "Permission denied"

```bash
# Vérifier les permissions
ls -la /opt/wedding-web/
chown -R wedding:wedding /opt/wedding-web/
```

### Certificat SSL expiré

```bash
# Renouveler manuellement
certbot renew
systemctl reload nginx
```

---

## 📊 Monitoring

### Vérifier l'espace disque

```bash
df -h
```

### Vérifier la mémoire

```bash
free -h
```

### Vérifier les processus

```bash
ps aux | grep wedding-web
```

---

## 🎉 C'est prêt !

Votre site de mariage est maintenant en ligne ! 🥳

- ✅ HTTPS configuré
- ✅ Backup automatique
- ✅ Service qui redémarre automatiquement
- ✅ Firewall configuré
- ✅ Données chiffrées

**Bon mariage ! 💍🎊**

---

## 📞 Support

En cas de problème, vérifiez :
1. Les logs : `/opt/wedding-web/logs/app.log`
2. Le statut du service : `systemctl status wedding-web`
3. La configuration Nginx : `nginx -t`
4. Le DNS : `dig votre-domaine.fr`

