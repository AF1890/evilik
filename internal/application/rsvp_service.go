package application

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"wedding-web/internal/domain"
	"wedding-web/internal/domain/ports"
)

var (
	ErrStorageFailure = errors.New("erreur de stockage")
)

// RSVPService gère la logique métier des RSVP
type RSVPService struct {
	storage ports.RSVPStorage
}

// NewRSVPService crée un nouveau service RSVP
func NewRSVPService(storage ports.RSVPStorage) *RSVPService {
	return &RSVPService{
		storage: storage,
	}
}

// SubmitRSVP enregistre un nouveau RSVP
func (s *RSVPService) SubmitRSVP(firstName, lastName string, willAttend bool, adultsCount, childrenCount int, allergies, message, ipAddress string) (*domain.RSVP, error) {
	// Création et validation
	rsvp, err := domain.NewRSVP(firstName, lastName, willAttend, adultsCount, childrenCount, allergies, message)
	if err != nil {
		return nil, err
	}

	// Génération d'un ID unique
	rsvp.ID = generateID()
	rsvp.IPAddress = ipAddress

	// Sauvegarde
	if err := s.storage.Save(rsvp); err != nil {
		return nil, ErrStorageFailure
	}

	return rsvp, nil
}

// ListRSVPs retourne tous les RSVP
func (s *RSVPService) ListRSVPs() ([]*domain.RSVP, error) {
	rsvps, err := s.storage.FindAll()
	if err != nil {
		return nil, ErrStorageFailure
	}
	return rsvps, nil
}

// DeleteRSVP supprime un RSVP par son ID
func (s *RSVPService) DeleteRSVP(id string) error {
	err := s.storage.Delete(id)
	if err != nil {
		return ErrStorageFailure
	}
	return nil
}

// GetRSVP retourne un RSVP par son ID
func (s *RSVPService) GetRSVP(id string) (*domain.RSVP, error) {
	rsvp, err := s.storage.FindByID(id)
	if err != nil {
		return nil, ErrStorageFailure
	}
	return rsvp, nil
}

// generateID génère un identifiant unique aléatoire
func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback sur timestamp en cas d'erreur (très rare)
		return hex.EncodeToString([]byte("fallback"))
	}
	return hex.EncodeToString(b)
}
