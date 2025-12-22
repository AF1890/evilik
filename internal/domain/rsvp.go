package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidRSVP      = errors.New("données RSVP invalides")
	ErrInvalidName      = errors.New("nom invalide")
	ErrInvalidGuests    = errors.New("nombre d'invités invalide")
	ErrMessageTooLong   = errors.New("message trop long")
	ErrAllergiesTooLong = errors.New("allergies trop longues")
)

// RSVP représente une réservation
type RSVP struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	WillAttend    bool      `json:"will_attend"`
	AdultsCount   int       `json:"adults_count"`
	ChildrenCount int       `json:"children_count"`
	Allergies     string    `json:"allergies"`
	Message       string    `json:"message"`
	SubmittedAt   time.Time `json:"submitted_at"`
	IPAddress     string    `json:"-"` // Ne pas persister l'IP
}

// NewRSVP crée une nouvelle réservation avec validation
func NewRSVP(firstName, lastName string, willAttend bool, adultsCount, childrenCount int, allergies, message string) (*RSVP, error) {
	// Validation du prénom
	firstName = strings.TrimSpace(firstName)
	if len(firstName) == 0 || len(firstName) > 100 {
		return nil, ErrInvalidName
	}

	// Validation du nom
	lastName = strings.TrimSpace(lastName)
	if len(lastName) == 0 || len(lastName) > 100 {
		return nil, ErrInvalidName
	}

	// Si absent, forcer les compteurs à 0
	if !willAttend {
		adultsCount = 0
		childrenCount = 0
		allergies = "" // Pas d'allergies pour les absents
	} else {
		// Validation du nombre d'invités (seulement si présent)
		if adultsCount < 0 || adultsCount > 20 {
			return nil, ErrInvalidGuests
		}
		if childrenCount < 0 || childrenCount > 20 {
			return nil, ErrInvalidGuests
		}
		if adultsCount == 0 && childrenCount == 0 {
			return nil, ErrInvalidGuests
		}
	}

	// Validation des allergies
	allergies = strings.TrimSpace(allergies)
	if len(allergies) > 500 {
		return nil, ErrAllergiesTooLong
	}

	// Validation du message
	message = strings.TrimSpace(message)
	if len(message) > 1000 {
		return nil, ErrMessageTooLong
	}

	return &RSVP{
		FirstName:     firstName,
		LastName:      lastName,
		WillAttend:    willAttend,
		AdultsCount:   adultsCount,
		ChildrenCount: childrenCount,
		Allergies:     allergies,
		Message:       message,
		SubmittedAt:   time.Now(),
	}, nil
}

// TotalGuests retourne le nombre total d'invités
func (r *RSVP) TotalGuests() int {
	return r.AdultsCount + r.ChildrenCount
}

// FullName retourne le nom complet
func (r *RSVP) FullName() string {
	return r.FirstName + " " + r.LastName
}
