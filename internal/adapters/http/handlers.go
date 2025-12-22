package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"wedding-web/internal/application"
)

// Handlers contient tous les handlers HTTP
type Handlers struct {
	rsvpService     *application.RSVPService
	planningService *application.PlanningService
	infoService     *application.InfoService
	calendarService *application.CalendarService
	exportService   *application.ExportService
	csrfManager     *CSRFManager
	templates       *template.Template
	templatesDir    string // Pour recharger les templates en dev
	isDev           bool   // Mode développement
	adminUsername   string // Nom d'utilisateur admin
	adminPassword   string // Mot de passe admin
}

// NewHandlers crée une nouvelle instance des handlers
func NewHandlers(
	rsvpService *application.RSVPService,
	planningService *application.PlanningService,
	infoService *application.InfoService,
	calendarService *application.CalendarService,
	csrfManager *CSRFManager,
	templatesDir string,
	isDev bool,
	adminUsername string,
	adminPassword string,
) (*Handlers, error) {
	// Charger les templates avec les partials
	tmpl, err := template.ParseGlob(filepath.Join(templatesDir, "partials", "*.html"))
	if err != nil {
		return nil, err
	}

	// Charger les templates de pages
	tmpl, err = tmpl.ParseFiles(
		filepath.Join(templatesDir, "home.html"),
		filepath.Join(templatesDir, "planning.html"),
		filepath.Join(templatesDir, "infos.html"),
		filepath.Join(templatesDir, "rsvp.html"),
		filepath.Join(templatesDir, "confirmation.html"),
		filepath.Join(templatesDir, "error.html"),
		filepath.Join(templatesDir, "admin.html"),
	)
	if err != nil {
		return nil, err
	}

	// Créer le service d'export
	exportService := application.NewExportService(rsvpService)

	return &Handlers{
		rsvpService:     rsvpService,
		planningService: planningService,
		infoService:     infoService,
		calendarService: calendarService,
		exportService:   exportService,
		csrfManager:     csrfManager,
		templates:       tmpl,
		templatesDir:    templatesDir,
		isDev:           isDev,
		adminUsername:   adminUsername,
		adminPassword:   adminPassword,
	}, nil
}

// reloadTemplates recharge les templates en mode dev
func (h *Handlers) reloadTemplates() {
	if !h.isDev {
		return
	}

	// Charger les partials d'abord
	tmpl, err := template.ParseGlob(filepath.Join(h.templatesDir, "partials", "*.html"))
	if err != nil {
		return
	}

	// Puis charger les templates de pages
	tmpl, err = tmpl.ParseFiles(
		filepath.Join(h.templatesDir, "home.html"),
		filepath.Join(h.templatesDir, "planning.html"),
		filepath.Join(h.templatesDir, "infos.html"),
		filepath.Join(h.templatesDir, "rsvp.html"),
		filepath.Join(h.templatesDir, "confirmation.html"),
		filepath.Join(h.templatesDir, "error.html"),
		filepath.Join(h.templatesDir, "admin.html"),
	)
	if err == nil {
		h.templates = tmpl
	}
}

// HomeHandler affiche la page d'accueil
func (h *Handlers) HomeHandler(w ResponseWriter, r *Request) error {
	h.reloadTemplates()

	data := map[string]interface{}{
		"Title": "A & G",
	}

	return h.templates.ExecuteTemplate(w, "home.html", data)
}

// PlanningHandler affiche le planning
func (h *Handlers) PlanningHandler(w ResponseWriter, r *Request) error {
	h.reloadTemplates()

	planning := h.planningService.GetPlanning()

	data := map[string]interface{}{
		"Title":    "Planning de la journée",
		"Planning": planning,
	}

	return h.templates.ExecuteTemplate(w, "planning.html", data)
}

// InfosHandler affiche les infos pratiques
func (h *Handlers) InfosHandler(w ResponseWriter, r *Request) error {
	h.reloadTemplates()

	info := h.infoService.GetPracticalInfo()

	data := map[string]interface{}{
		"Title": "Informations pratiques",
		"Info":  info,
	}

	return h.templates.ExecuteTemplate(w, "infos.html", data)
}

// RSVPGetHandler affiche le formulaire RSVP
func (h *Handlers) RSVPGetHandler(w ResponseWriter, r *Request) error {
	h.reloadTemplates()

	// Récupérer ou créer une session
	sessionID := getOrCreateSession(w, r.Request)

	// Générer un token CSRF
	csrfToken, err := h.csrfManager.GenerateToken(sessionID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Title":     "Confirmez votre présence",
		"CSRFToken": csrfToken,
	}

	return h.templates.ExecuteTemplate(w, "rsvp.html", data)
}

// RSVPPostHandler traite la soumission du formulaire RSVP
func (h *Handlers) RSVPPostHandler(w ResponseWriter, r *Request) error {
	// Vérifier le Content-Type
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content-Type invalide"))
		return nil
	}

	// Parser le formulaire
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Formulaire invalide"))
		return nil
	}

	// Vérifier le token CSRF
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Session invalide"))
		return nil
	}

	csrfToken := r.FormValue("csrf_token")
	if !h.csrfManager.ValidateToken(cookie.Value, csrfToken) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Token CSRF invalide"))
		return nil
	}

	// Vérifier le honeypot (champ caché anti-spam)
	honeypot := r.FormValue("website")
	if honeypot != "" {
		// C'est probablement un bot
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Requête invalide"))
		return nil
	}

	// Récupérer les données
	// Prénom/Nom peuvent venir de deux sections différentes
	firstName := strings.TrimSpace(r.FormValue("presence_first_name"))
	if firstName == "" {
		firstName = strings.TrimSpace(r.FormValue("absence_first_name"))
	}

	lastName := strings.TrimSpace(r.FormValue("presence_last_name"))
	if lastName == "" {
		lastName = strings.TrimSpace(r.FormValue("absence_last_name"))
	}

	willAttendStr := r.FormValue("will_attend")
	willAttend := willAttendStr == "yes"

	adultsCount, _ := strconv.Atoi(r.FormValue("adults_count"))
	childrenCount, _ := strconv.Atoi(r.FormValue("children_count"))
	allergies := strings.TrimSpace(r.FormValue("allergies"))

	// Message peut venir de deux champs différents (présence ou absence)
	message := strings.TrimSpace(r.FormValue("presence_message"))
	if message == "" {
		message = strings.TrimSpace(r.FormValue("absence_message"))
	}

	// Récupérer l'IP
	ip := getClientIP(r.Request)

	// Soumettre le RSVP
	rsvp, err := h.rsvpService.SubmitRSVP(firstName, lastName, willAttend, adultsCount, childrenCount, allergies, message, ip)
	if err != nil {
		data := map[string]interface{}{
			"Title": "Erreur",
			"Error": err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		return h.templates.ExecuteTemplate(w, "error.html", data)
	}

	// Redirection vers la page de confirmation
	data := map[string]interface{}{
		"Title":     "Confirmation",
		"FirstName": rsvp.FirstName,
		"LastName":  rsvp.LastName,
		"Total":     rsvp.TotalGuests(),
	}

	return h.templates.ExecuteTemplate(w, "confirmation.html", data)
}

// CalendarHandler génère et retourne un fichier .ics
func (h *Handlers) CalendarHandler(w ResponseWriter, r *Request) error {
	planning := h.planningService.GetPlanning()

	// Générer le fichier ICS
	icsData, err := h.calendarService.GenerateICS(planning)
	if err != nil {
		return err
	}

	// Définir les headers pour le téléchargement
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=mariage-2026.ics")
	w.Header().Set("Cache-Control", "no-cache")

	w.Write(icsData)
	return nil
}

// HealthHandler retourne le statut de santé du service
func (h *Handlers) HealthHandler(w ResponseWriter, r *Request) error {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
	return nil
}

// NotFoundHandler gère les 404
func (h *Handlers) NotFoundHandler(w ResponseWriter, r *Request) error {
	w.WriteHeader(http.StatusNotFound)

	data := map[string]interface{}{
		"Title":   "Page non trouvée",
		"Message": "La page que vous cherchez n'existe pas.",
	}

	return h.templates.ExecuteTemplate(w, "error.html", data)
}

// AdminHandler affiche la liste des RSVP (protégé par mot de passe)
func (h *Handlers) AdminHandler(w ResponseWriter, r *Request) error {
	// Vérifier si l'admin est configuré
	if h.adminPassword == "" || h.adminUsername == "" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	// Vérifier l'authentification
	username, password, ok := r.BasicAuth()
	if !ok || username != h.adminUsername || password != h.adminPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="Administration - RSVP Mariage"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authentification requise"))
		return nil
	}

	h.reloadTemplates()

	// Récupérer tous les RSVP
	rsvps, err := h.rsvpService.ListRSVPs()
	if err != nil {
		return err
	}

	// Calculer les statistiques
	totalPersonnes := 0
	totalAdultes := 0
	totalEnfants := 0
	totalPresents := 0
	totalAbsents := 0

	for _, rsvp := range rsvps {
		if rsvp.WillAttend {
			totalPresents++
			totalAdultes += rsvp.AdultsCount
			totalEnfants += rsvp.ChildrenCount
			totalPersonnes += rsvp.TotalGuests()
		} else {
			totalAbsents++
		}
	}

	data := map[string]interface{}{
		"Title":          "Administration - RSVPs",
		"RSVPs":          rsvps,
		"TotalRSVPs":     len(rsvps),
		"TotalPresents":  totalPresents,
		"TotalAbsents":   totalAbsents,
		"TotalPersonnes": totalPersonnes,
		"TotalAdultes":   totalAdultes,
		"TotalEnfants":   totalEnfants,
	}

	return h.templates.ExecuteTemplate(w, "admin.html", data)
}

// AdminExportHandler génère et télécharge un fichier Excel des RSVP
func (h *Handlers) AdminExportHandler(w ResponseWriter, r *Request) error {
	// Vérifier si l'admin est configuré
	if h.adminPassword == "" || h.adminUsername == "" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	// Vérifier l'authentification
	username, password, ok := r.BasicAuth()
	if !ok || username != h.adminUsername || password != h.adminPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="Administration - RSVP Mariage"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authentification requise"))
		return nil
	}

	// Générer le fichier Excel
	file, err := h.exportService.ExportRSVPsToExcel()
	if err != nil {
		return err
	}
	defer file.Close()

	// Nom du fichier
	filename := h.exportService.GetFileName()

	// Headers pour le téléchargement
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Écrire le fichier dans la réponse
	return file.Write(w)
}

// AdminDeleteHandler supprime un RSVP
func (h *Handlers) AdminDeleteHandler(w ResponseWriter, r *Request) error {
	// Vérifier si l'admin est configuré
	if h.adminPassword == "" || h.adminUsername == "" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	// Vérifier l'authentification
	username, password, ok := r.BasicAuth()
	if !ok || username != h.adminUsername || password != h.adminPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="Administration - RSVP Mariage"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authentification requise"))
		return nil
	}

	// Récupérer l'ID depuis l'URL
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID manquant"))
		return nil
	}

	// Supprimer le RSVP
	err := h.rsvpService.DeleteRSVP(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erreur lors de la suppression"))
		return nil
	}

	// Rediriger vers la page admin
	http.Redirect(w, r.Request, "/admin", http.StatusSeeOther)
	return nil
}
