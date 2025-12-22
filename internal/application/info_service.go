package application

import "wedding-web/internal/domain"

// InfoService gère la logique métier des infos pratiques
type InfoService struct {
	info *domain.PracticalInfo
}

// NewInfoService crée un nouveau service d'infos
func NewInfoService() *InfoService {
	return &InfoService{
		info: domain.GetDefaultPracticalInfo(),
	}
}

// GetPracticalInfo retourne les informations pratiques
func (s *InfoService) GetPracticalInfo() *domain.PracticalInfo {
	return s.info
}
