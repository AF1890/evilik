package i18n

import (
	"time"
)

// Lang représente une langue
type Lang string

const (
	FR Lang = "fr"
	DE Lang = "de"
)

// Translations contient toutes les traductions
type Translations struct {
	lang Lang
	data map[string]string
}

// NewTranslations crée une nouvelle instance de traductions
func NewTranslations(lang Lang) *Translations {
	if lang != FR && lang != DE {
		lang = FR // Par défaut français
	}
	return &Translations{
		lang: lang,
		data: getTranslations(lang),
	}
}

// T retourne une traduction
func (t *Translations) T(key string) string {
	if val, ok := t.data[key]; ok {
		return val
	}
	return key // Retourne la clé si traduction non trouvée
}

// Lang retourne la langue courante
func (t *Translations) Lang() string {
	return string(t.lang)
}

// FormatDate formate une date selon la langue
func (t *Translations) FormatDate(date time.Time) string {
	if t.lang == DE {
		return date.Format("2. January 2006")
	}
	return date.Format("2 janvier 2006")
}

// FormatDateTime formate une date/heure selon la langue
func (t *Translations) FormatDateTime(date time.Time) string {
	if t.lang == DE {
		return date.Format("2. January 2006 um 15:04 Uhr")
	}
	return date.Format("2 janvier 2006 à 15h04")
}

// getTranslations retourne toutes les traductions pour une langue
func getTranslations(lang Lang) map[string]string {
	switch lang {
	case DE:
		return germanTranslations
	default:
		return frenchTranslations
	}
}

// frenchTranslations - Traductions françaises
var frenchTranslations = map[string]string{
	// Navigation
	"nav.home":     "Accueil",
	"nav.planning": "Planning",
	"nav.info":     "Infos pratiques",
	"nav.rsvp":     "RSVP",

	// Page d'accueil
	"home.title":       "Aylin et Guillaume",
	"home.date":        "11 juillet 2026 à Cély-en-Bière",
	"home.intro":       "16 ans de vie commune, deux enfants... et maintenant,\non embarque pour une nouvelle aventure: on se marie !\nHâte de vous avoir avec nous !",
	"home.cta_button":  "Confirmer ma présence",
	"home.card1_title": "Planning",
	"home.card1_desc":  "Découvrez le déroulement de la journée",
	"home.card2_title": "Infos pratiques",
	"home.card2_desc":  "Lieu, accès, hébergement...",
	"home.card3_title": "Confirmer votre présence",
	"home.card3_desc":  "Merci de nous répondre avant le 1er mars 2026",

	// Planning
	"planning.title":          "Planning de la journée",
	"planning.subtitle":       "Déroulement du mariage",
	"planning.download":       "Télécharger au format .ics",
	"planning.ceremony_title": "Cérémonie civile",
	"planning.ceremony_desc":  "Notre union officielle à la mairie",
	"planning.cocktail_title": "Cocktail & Vin d'honneur",
	"planning.cocktail_desc":  "Moments de convivialité et de partage",
	"planning.photo_title":    "Séance photo",
	"planning.photo_desc":     "Laissez-nous quelques instants pour immortaliser ce jour",
	"planning.dinner_title":   "Dîner & Soirée",
	"planning.dinner_desc":    "Repas et fête jusqu'au bout de la nuit !",
	"planning.info_box":       "Note importante : Les horaires peuvent légèrement varier.",

	// Infos pratiques
	"info.title":               "Informations pratiques",
	"info.venue_title":         "Le lieu",
	"info.venue_content":       "Le mariage civil aura lieu à la mairie de Cély.\nLa cérémonie laïque se tiendra ensuite chez nous, au 8 rue du Bois Beaudoin, à Cély-en-Bière.\nLa soirée se poursuivra à la Bergerie de Villiers-en-Bière.\nAdresse : rue de la Bascule, 77190 Villiers-en-Bière.",
	"info.venue_name":          "Château de Cély",
	"info.venue_address":       "1 Rue du Château, 77930 Cély-en-Bière",
	"info.venue_desc":          "Un magnifique château du XVIIIe siècle situé en pleine nature, à 50 minutes de Paris.",
	"info.venue_map":           "Voir sur la carte",
	"info.access_title":        "Accès",
	"info.access_desc":         "En voiture : Autoroute A6, sortie Fontainebleau\nEn train : Gare de Fontainebleau-Avon + navette (nous contacter)",
	"info.accommodation_title": "Hébergement",
	"info.accommodation_desc":  "Plusieurs hôtels à proximité :\n• Hôtel de Londres (Fontainebleau)\n• Hôtel Belle Fontainebleau\n• Chambres d'hôtes locales\nN'hésitez pas à nous contacter pour des recommandations.",
	"info.dresscode_title":     "Tenue",
	"info.dresscode_desc":      "Tenue de cérémonie souhaitée. Chic et élégant !",
	"info.contact_title":       "Contact",
	"info.contact_desc":        "Pour toute question :\naylin@exemple.com\nguillaume@exemple.com",

	// RSVP
	"rsvp.title":                 "Confirmez votre présence",
	"rsvp.subtitle":              "Merci de répondre avant le 1er mars 2026",
	"rsvp.attendance":            "Serez-vous présent(e) ?",
	"rsvp.attendance_yes":        "Oui, je serai là ! 🎉",
	"rsvp.attendance_no":         "Non, je ne pourrai pas 😢",
	"rsvp.firstname":             "Prénom",
	"rsvp.lastname":              "Nom",
	"rsvp.adults":                "Nombre d'adultes",
	"rsvp.children":              "Nombre d'enfants",
	"rsvp.allergies":             "Allergies / Régimes alimentaires",
	"rsvp.allergies_placeholder": "Végétarien, sans gluten, etc.",
	"rsvp.message":               "Un petit mot pour nous ?",
	"rsvp.message_placeholder":   "Partagez votre joie avec nous...",
	"rsvp.message_absence":       "Message (optionnel)",
	"rsvp.message_absence_ph":    "Nous espérons vous voir une prochaine fois...",
	"rsvp.submit":                "Envoyer ma réponse",
	"rsvp.privacy_note":          "Les informations collectées sont uniquement utilisées pour l'organisation du mariage et ne seront pas partagées avec des tiers.",
	"rsvp.confirmation":          "Merci pour votre réponse !",
	"rsvp.confirmation_text":     "Nous avons bien reçu votre confirmation. À très bientôt !",
	"rsvp.back":                  "Retour à l'accueil",

	// Footer
	"footer.copyright": "© 2026 Aylin & Guillaume",

	// Erreurs
	"error.title":              "Une erreur est survenue",
	"error.desc":               "Désolé, quelque chose s'est mal passé.",
	"error.back":               "Retour à l'accueil",
	"error.contact":            "Si le problème persiste, contactez-nous.",
	"error.invalid_name":       "Le prénom et le nom sont obligatoires (maximum 100 caractères)",
	"error.invalid_guests":     "Le nombre d'invités est invalide (au moins 1 adulte ou enfant requis)",
	"error.message_too_long":   "Le message est trop long (maximum 1000 caractères)",
	"error.allergies_too_long": "Les allergies sont trop longues (maximum 500 caractères)",
}

// germanTranslations - Deutsche Übersetzungen
var germanTranslations = map[string]string{
	// Navigation
	"nav.home":     "Startseite",
	"nav.planning": "Tagesablauf",
	"nav.info":     "Praktische Infos",
	"nav.rsvp":     "Zusagen",

	// Startseite
	"home.title":       "Aylin und Guillaume",
	"home.date":        "11. Juli 2026 in Cély-en-Bière",
	"home.intro":       "16 Jahre gemeinsames Leben, zwei Kinder... und jetzt\nstarten wir ein neues Abenteuer: Wir heiraten!\nWir freuen uns auf euch!",
	"home.cta_button":  "Zusage bestätigen",
	"home.card1_title": "Tagesablauf",
	"home.card1_desc":  "Entdecken Sie den Ablauf des Tages",
	"home.card2_title": "Praktische Infos",
	"home.card2_desc":  "Ort, Anfahrt, Unterkunft...",
	"home.card3_title": "Zusage bestätigen",
	"home.card3_desc":  "Bitte antworten Sie uns bis zum 1. März 2026",

	// Tagesablauf
	"planning.title":          "Tagesablauf",
	"planning.subtitle":       "Ablauf der Hochzeit",
	"planning.download":       "Als .ics herunterladen",
	"planning.ceremony_title": "Standesamtliche Trauung",
	"planning.ceremony_desc":  "Unsere offizielle Trauung im Rathaus",
	"planning.cocktail_title": "Sektempfang & Ehrenwein",
	"planning.cocktail_desc":  "Momente der Geselligkeit und des Teilens",
	"planning.photo_title":    "Fotoshooting",
	"planning.photo_desc":     "Gönnen Sie uns ein paar Augenblicke, um diesen Tag festzuhalten",
	"planning.dinner_title":   "Abendessen & Party",
	"planning.dinner_desc":    "Essen und Feiern bis in die Nacht!",
	"planning.info_box":       "Wichtiger Hinweis: Die Zeiten können sich leicht ändern.",

	// Praktische Infos
	"info.title":               "Praktische Informationen",
	"info.venue_title":         "Der Ort",
	"info.venue_content":       "Die standesamtliche Trauung findet im Rathaus von Cély statt.\nDie weltliche Zeremonie findet anschließend bei uns statt, 8 rue du Bois Beaudoin, in Cély-en-Bière.\nDer Abend wird in der Bergerie de Villiers-en-Bière fortgesetzt.\nAdresse: rue de la Bascule, 77190 Villiers-en-Bière.",
	"info.venue_name":          "Schloss Cély",
	"info.venue_address":       "1 Rue du Château, 77930 Cély-en-Bière",
	"info.venue_desc":          "Ein wunderschönes Schloss aus dem 18. Jahrhundert mitten in der Natur, 50 Minuten von Paris entfernt.",
	"info.venue_map":           "Auf der Karte ansehen",
	"info.access_title":        "Anfahrt",
	"info.access_desc":         "Mit dem Auto: Autobahn A6, Ausfahrt Fontainebleau\nMit dem Zug: Bahnhof Fontainebleau-Avon + Shuttle (kontaktieren Sie uns)",
	"info.accommodation_title": "Unterkunft",
	"info.accommodation_desc":  "Mehrere Hotels in der Nähe:\n• Hôtel de Londres (Fontainebleau)\n• Hôtel Belle Fontainebleau\n• Lokale Gästehäuser\nKontaktieren Sie uns gerne für Empfehlungen.",
	"info.dresscode_title":     "Kleiderordnung",
	"info.dresscode_desc":      "Festliche Kleidung erwünscht. Schick und elegant!",
	"info.contact_title":       "Kontakt",
	"info.contact_desc":        "Für alle Fragen:\naylin@beispiel.com\nguillaume@beispiel.com",

	// RSVP / Zusage
	"rsvp.title":                 "Bestätigen Sie Ihre Anwesenheit",
	"rsvp.subtitle":              "Bitte antworten Sie bis zum 1. März 2026",
	"rsvp.attendance":            "Werden Sie dabei sein?",
	"rsvp.attendance_yes":        "Ja, ich werde da sein! 🎉",
	"rsvp.attendance_no":         "Nein, ich kann leider nicht 😢",
	"rsvp.firstname":             "Vorname",
	"rsvp.lastname":              "Nachname",
	"rsvp.adults":                "Anzahl Erwachsene",
	"rsvp.children":              "Anzahl Kinder",
	"rsvp.allergies":             "Allergien / Ernährungsweise",
	"rsvp.allergies_placeholder": "Vegetarisch, glutenfrei, usw.",
	"rsvp.message":               "Eine kleine Nachricht für uns?",
	"rsvp.message_placeholder":   "Teilen Sie Ihre Freude mit uns...",
	"rsvp.message_absence":       "Nachricht (optional)",
	"rsvp.message_absence_ph":    "Wir hoffen, Sie bald zu sehen...",
	"rsvp.submit":                "Antwort senden",
	"rsvp.privacy_note":          "Die gesammelten Informationen werden ausschließlich für die Organisation der Hochzeit verwendet und nicht an Dritte weitergegeben.",
	"rsvp.confirmation":          "Vielen Dank für Ihre Antwort!",
	"rsvp.confirmation_text":     "Wir haben Ihre Bestätigung erhalten. Bis bald!",
	"rsvp.back":                  "Zurück zur Startseite",

	// Footer
	"footer.copyright": "© 2026 Aylin & Guillaume",

	// Fehler
	"error.title":              "Ein Fehler ist aufgetreten",
	"error.desc":               "Entschuldigung, etwas ist schief gelaufen.",
	"error.back":               "Zurück zur Startseite",
	"error.contact":            "Wenn das Problem weiterhin besteht, kontaktieren Sie uns bitte.",
	"error.invalid_name":       "Vorname und Nachname sind erforderlich (maximal 100 Zeichen)",
	"error.invalid_guests":     "Die Anzahl der Gäste ist ungültig (mindestens 1 Erwachsener oder Kind erforderlich)",
	"error.message_too_long":   "Die Nachricht ist zu lang (maximal 1000 Zeichen)",
	"error.allergies_too_long": "Die Allergieinformationen sind zu lang (maximal 500 Zeichen)",
}
