package application

import (
	"testing"
	"wedding-web/internal/domain"
)

// Mock storage pour les tests
type mockStorage struct {
	rsvps []*domain.RSVP
	err   error
}

func (m *mockStorage) Save(rsvp *domain.RSVP) error {
	if m.err != nil {
		return m.err
	}
	m.rsvps = append(m.rsvps, rsvp)
	return nil
}

func (m *mockStorage) FindAll() ([]*domain.RSVP, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rsvps, nil
}

func (m *mockStorage) FindByID(id string) (*domain.RSVP, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, rsvp := range m.rsvps {
		if rsvp.ID == id {
			return rsvp, nil
		}
	}
	return nil, domain.ErrInvalidRSVP
}

func (m *mockStorage) Delete(id string) error {
	if m.err != nil {
		return m.err
	}
	for i, rsvp := range m.rsvps {
		if rsvp.ID == id {
			m.rsvps = append(m.rsvps[:i], m.rsvps[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestRSVPService_SubmitRSVP(t *testing.T) {
	storage := &mockStorage{rsvps: []*domain.RSVP{}}
	service := NewRSVPService(storage)

	rsvp, err := service.SubmitRSVP("Jean", "Dupont", true, 2, 1, "Aucune", "Message", "127.0.0.1")

	if err != nil {
		t.Fatalf("SubmitRSVP() error = %v", err)
	}

	if rsvp == nil {
		t.Fatal("SubmitRSVP() returned nil")
	}

	if rsvp.ID == "" {
		t.Error("RSVP ID is empty")
	}

	if rsvp.FirstName != "Jean" {
		t.Errorf("FirstName = %s, want Jean", rsvp.FirstName)
	}

	// Vérifier que le RSVP a été sauvegardé
	if len(storage.rsvps) != 1 {
		t.Errorf("Storage contains %d RSVPs, want 1", len(storage.rsvps))
	}
}

func TestRSVPService_SubmitRSVP_Invalid(t *testing.T) {
	storage := &mockStorage{rsvps: []*domain.RSVP{}}
	service := NewRSVPService(storage)

	// Test avec des données invalides
	_, err := service.SubmitRSVP("", "Dupont", true, 1, 0, "", "", "127.0.0.1")

	if err == nil {
		t.Error("SubmitRSVP() should return an error for invalid data")
	}
}

func TestRSVPService_ListRSVPs(t *testing.T) {
	storage := &mockStorage{rsvps: []*domain.RSVP{}}
	service := NewRSVPService(storage)

	// Ajouter quelques RSVPs
	service.SubmitRSVP("Jean", "Dupont", true, 2, 0, "", "", "127.0.0.1")
	service.SubmitRSVP("Marie", "Martin", true, 1, 1, "", "", "127.0.0.1")

	rsvps, err := service.ListRSVPs()

	if err != nil {
		t.Fatalf("ListRSVPs() error = %v", err)
	}

	if len(rsvps) != 2 {
		t.Errorf("ListRSVPs() returned %d RSVPs, want 2", len(rsvps))
	}
}

func TestCalendarService_GenerateICS(t *testing.T) {
	service := NewCalendarService()
	planning := domain.GetDefaultPlanning()

	icsData, err := service.GenerateICS(planning)

	if err != nil {
		t.Fatalf("GenerateICS() error = %v", err)
	}

	if len(icsData) == 0 {
		t.Error("GenerateICS() returned empty data")
	}

	// Vérifier que le fichier contient les marqueurs ICS
	icsString := string(icsData)

	if !contains(icsString, "BEGIN:VCALENDAR") {
		t.Error("ICS data missing BEGIN:VCALENDAR")
	}

	if !contains(icsString, "END:VCALENDAR") {
		t.Error("ICS data missing END:VCALENDAR")
	}

	if !contains(icsString, "BEGIN:VEVENT") {
		t.Error("ICS data missing BEGIN:VEVENT")
	}
}

func TestPlanningService(t *testing.T) {
	service := NewPlanningService()
	planning := service.GetPlanning()

	if planning == nil {
		t.Fatal("GetPlanning() returned nil")
	}

	if len(planning.Events) == 0 {
		t.Error("Planning has no events")
	}
}

func TestInfoService(t *testing.T) {
	service := NewInfoService()
	info := service.GetPracticalInfo()

	if info == nil {
		t.Fatal("GetPracticalInfo() returned nil")
	}

	if info.Venue.Title == "" {
		t.Error("Venue title is empty")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && (s == substr || len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && s[len(s)-len(substr):] == substr || len(s) > len(substr) && containsInMiddle(s, substr))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
