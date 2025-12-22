package application

import "wedding-web/internal/domain"

// PlanningService gère la logique métier du planning
type PlanningService struct {
	planning *domain.Planning
}

// NewPlanningService crée un nouveau service de planning
func NewPlanningService() *PlanningService {
	return &PlanningService{
		planning: domain.GetDefaultPlanning(),
	}
}

// GetPlanning retourne le planning complet
func (s *PlanningService) GetPlanning() *domain.Planning {
	return s.planning
}
