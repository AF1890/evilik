package ports

import "wedding-web/internal/domain"

// CalendarGenerator génère des fichiers calendar (.ics)
type CalendarGenerator interface {
	GenerateICS(planning *domain.Planning) ([]byte, error)
}
