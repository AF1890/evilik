package domain

import (
	"testing"
)

func TestNewRSVP(t *testing.T) {
	tests := []struct {
		name          string
		firstName     string
		lastName      string
		adultsCount   int
		childrenCount int
		allergies     string
		message       string
		wantErr       error
	}{
		{
			name:          "Valid RSVP",
			firstName:     "Jean",
			lastName:      "Dupont",
			adultsCount:   2,
			childrenCount: 1,
			allergies:     "Aucune",
			message:       "Hâte d'être là !",
			wantErr:       nil,
		},
		{
			name:          "Empty first name",
			firstName:     "",
			lastName:      "Dupont",
			adultsCount:   1,
			childrenCount: 0,
			wantErr:       ErrInvalidName,
		},
		{
			name:          "Empty last name",
			firstName:     "Jean",
			lastName:      "",
			adultsCount:   1,
			childrenCount: 0,
			wantErr:       ErrInvalidName,
		},
		{
			name:          "Name too long",
			firstName:     string(make([]byte, 101)),
			lastName:      "Dupont",
			adultsCount:   1,
			childrenCount: 0,
			wantErr:       ErrInvalidName,
		},
		{
			name:          "No guests",
			firstName:     "Jean",
			lastName:      "Dupont",
			adultsCount:   0,
			childrenCount: 0,
			wantErr:       ErrInvalidGuests,
		},
		{
			name:          "Negative adults",
			firstName:     "Jean",
			lastName:      "Dupont",
			adultsCount:   -1,
			childrenCount: 0,
			wantErr:       ErrInvalidGuests,
		},
		{
			name:          "Too many adults",
			firstName:     "Jean",
			lastName:      "Dupont",
			adultsCount:   25,
			childrenCount: 0,
			wantErr:       ErrInvalidGuests,
		},
		{
			name:          "Allergies too long",
			firstName:     "Jean",
			lastName:      "Dupont",
			adultsCount:   1,
			childrenCount: 0,
			allergies:     string(make([]byte, 501)),
			wantErr:       ErrAllergiesTooLong,
		},
		{
			name:          "Message too long",
			firstName:     "Jean",
			lastName:      "Dupont",
			adultsCount:   1,
			childrenCount: 0,
			message:       string(make([]byte, 1001)),
			wantErr:       ErrMessageTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsvp, err := NewRSVP(tt.firstName, tt.lastName, true, tt.adultsCount, tt.childrenCount, tt.allergies, tt.message)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewRSVP() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewRSVP() unexpected error = %v", err)
				return
			}

			if rsvp == nil {
				t.Error("NewRSVP() returned nil")
				return
			}

			if rsvp.FirstName != tt.firstName {
				t.Errorf("FirstName = %v, want %v", rsvp.FirstName, tt.firstName)
			}
			if rsvp.LastName != tt.lastName {
				t.Errorf("LastName = %v, want %v", rsvp.LastName, tt.lastName)
			}
		})
	}
}

func TestRSVP_TotalGuests(t *testing.T) {
	rsvp, _ := NewRSVP("Jean", "Dupont", true, 2, 3, "", "")

	total := rsvp.TotalGuests()
	expected := 5

	if total != expected {
		t.Errorf("TotalGuests() = %d, want %d", total, expected)
	}
}

func TestRSVP_FullName(t *testing.T) {
	rsvp, _ := NewRSVP("Jean", "Dupont", true, 1, 0, "", "")

	fullName := rsvp.FullName()
	expected := "Jean Dupont"

	if fullName != expected {
		t.Errorf("FullName() = %s, want %s", fullName, expected)
	}
}

func TestGetDefaultPlanning(t *testing.T) {
	planning := GetDefaultPlanning()

	if planning == nil {
		t.Fatal("GetDefaultPlanning() returned nil")
	}

	if len(planning.Events) == 0 {
		t.Error("GetDefaultPlanning() returned no events")
	}

	if planning.WeddingDate.IsZero() {
		t.Error("WeddingDate is zero")
	}

	// Vérifier que la date est bien le 11 juillet 2026
	if planning.WeddingDate.Year() != 2026 || planning.WeddingDate.Month() != 7 || planning.WeddingDate.Day() != 11 {
		t.Errorf("WeddingDate = %v, want 2026-07-11", planning.WeddingDate)
	}
}

func TestGetDefaultPracticalInfo(t *testing.T) {
	info := GetDefaultPracticalInfo()

	if info == nil {
		t.Fatal("GetDefaultPracticalInfo() returned nil")
	}

	if info.Venue.Title == "" {
		t.Error("Venue title is empty")
	}
}
