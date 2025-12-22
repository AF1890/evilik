package domain

import "time"

// PlanningEvent représente un événement du planning
type PlanningEvent struct {
	Title        string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	Location     string
	Address      string
	HideTime     bool // Masquer l'heure (par défaut: affiché)
	HideLocation bool // Masquer la localisation (par défaut: affiché)
}

// Planning représente le planning complet de la journée
type Planning struct {
	WeddingDate time.Time
	Events      []PlanningEvent
}

// GetDefaultPlanning retourne le planning par défaut
func GetDefaultPlanning() *Planning {
	weddingDate := time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC)

	return &Planning{
		WeddingDate: weddingDate,
		Events: []PlanningEvent{
			{
				Title:       "Cérémonie civile",
				StartTime:   time.Date(2026, 7, 11, 14, 00, 0, 0, time.UTC),
				EndTime:     time.Date(2026, 7, 11, 15, 00, 0, 0, time.UTC),
				Location:    "Mairie",
				Address:     "13 Rue de la Mairie, 77930 Cély",
				Description: "Se garer dans le parking de la mairie",
			},
			{
				Title:     "Cérémonie laïque",
				StartTime: time.Date(2026, 7, 11, 15, 30, 0, 0, time.UTC),
				EndTime:   time.Date(2026, 7, 11, 16, 30, 0, 0, time.UTC),
				Location:  "Chez nous",
				Address:   "8 rue du bois beaudoin, 77930 Cély",
			},
			{
				Title:        "Séance photo",
				Description:  "Photos des mariés et des invités",
				HideTime:     true,
				HideLocation: true,
			},
			{
				Title:       "Vin d'honneur",
				StartTime:   time.Date(2026, 7, 11, 18, 00, 0, 0, time.UTC),
				EndTime:     time.Date(2026, 7, 11, 20, 00, 0, 0, time.UTC),
				Location:    "La Bergerie",
				Address:     "Rue de la Bascule, 77190 Villiers-en-Bière",
				Description: "Un grand parking est disponible sur place.",
			},
			{
				Title:        "Diner",
				HideTime:     true,
				HideLocation: true,
			},
		},
	}
}
