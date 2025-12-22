package ports

import "wedding-web/internal/domain"

// RSVPStorage définit le port pour la persistance des RSVP
type RSVPStorage interface {
	Save(rsvp *domain.RSVP) error
	FindAll() ([]*domain.RSVP, error)
	FindByID(id string) (*domain.RSVP, error)
	Delete(id string) error
}
