# 📊 Page Administration - Guide d'utilisation

## Configuration

### Étape 1 : Définir les identifiants admin

Dans votre fichier `.env`, ajoutez :

```bash
ADMIN_USERNAME=votre-nom-utilisateur
ADMIN_PASSWORD=votre-mot-de-passe-secret
```

**Exemple pour le développement** :
```bash
ADMIN_USERNAME=admin
ADMIN_PASSWORD=MonMotDePasseSecurise2026!
```

**Exemple pour la production** :
```bash
ADMIN_USERNAME=aylin
ADMIN_PASSWORD=SuperMotDePasseTresComplique2026!#Secure
```

### Étape 2 : Redémarrer le serveur

```bash
make watch
```

Vous verrez dans les logs :
```
✅ Page admin activée sur /admin (utilisateur: admin)
```

---

## Accès à la page admin

### En développement local

1. **URL** : http://:8080/admin
2. **Authentification** :
   - Utilisateur : celui défini dans `ADMIN_USERNAME`
   - Mot de passe : celui défini dans `ADMIN_PASSWORD`

### En production

1. **URL** : https://votre-domaine.com/admin
2. **Authentification** : même identifiants que configurés

---

## Fonctionnalités

### Statistiques en temps réel
- ✅ Nombre total de confirmations
- ✅ Nombre total de personnes (adultes + enfants)
- ✅ Répartition adultes/enfants

### Liste des RSVP
Chaque confirmation affiche :
- 👤 Nom complet
- 📅 Date et heure de soumission
- 👥 Nombre d'adultes
- 👶 Nombre d'enfants
- 📊 Total de personnes
- 🍽️ Allergies/régimes (si renseignés)
- 💬 Message personnel (si renseigné)

---

## Sécurité

✅ **Protégé par mot de passe** : Seuls vous pouvez y accéder  
✅ **Données chiffrées** : Les RSVP sont stockés de manière sécurisée  
✅ **Pas d'email requis** : Aucune configuration SMTP nécessaire  
✅ **Consultation à tout moment** : Vérifiez les confirmations quand vous voulez

---

## ⚙️ Configuration : Dev vs Production

### En Développement (ENV=dev)

```bash
ENV=dev
ADMIN_USERNAME=admin
ADMIN_PASSWORD=test123
ENABLE_HSTS=false
```

- Mot de passe simple OK
- HTTP acceptable (localhost)
- Templates rechargés à chaque requête

### En Production (ENV=prod)

```bash
ENV=prod
ADMIN_USERNAME=aylin
ADMIN_PASSWORD=SuperMotDePasseTresComplique2026!#Secure
ENABLE_HSTS=true
```

**Règles de sécurité en production** :
1. ✅ **Mot de passe fort** : Au moins 16 caractères, lettres, chiffres, symboles
2. ✅ **HTTPS obligatoire** : Ne jamais utiliser HTTP en prod
3. ✅ **Username personnalisé** : Évitez "admin", "root", etc.
4. ✅ **HSTS activé** : Force HTTPS dans le navigateur
5. ✅ **Ne partagez jamais** les identifiants

**URL en production** : `https://votre-domaine.com/admin`

---

## Conseils

💡 **Gardez l'onglet ouvert** : Rafraîchissez (F5) pour voir les nouvelles confirmations  
💡 **Export manuel** : Faites des captures d'écran si besoin  
💡 **Vérifiez régulièrement** : Consultez la page 1-2 fois par semaine

---

C'est tout ! Beaucoup plus simple que les emails 😊

