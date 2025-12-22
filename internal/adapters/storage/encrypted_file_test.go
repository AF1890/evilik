package storage

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
	"wedding-web/internal/domain"
)

func TestEncryptedFileStorage(t *testing.T) {
	// Créer un répertoire temporaire
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_rsvp.json")

	// Générer une clé de test (32 bytes)
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	encodedKey := base64.StdEncoding.EncodeToString(key)

	// Créer le storage
	storage, err := NewEncryptedFileStorage(filePath, encodedKey)
	if err != nil {
		t.Fatalf("Erreur création storage: %v", err)
	}

	// Test 1: Sauvegarder un RSVP
	rsvp1, err := domain.NewRSVP("Jean", "Dupont", true, 2, 1, "Aucune", "Hâte d'être là !")
	if err != nil {
		t.Fatalf("Erreur création RSVP: %v", err)
	}
	rsvp1.ID = "test-id-1"

	if err := storage.Save(rsvp1); err != nil {
		t.Fatalf("Erreur sauvegarde RSVP: %v", err)
	}

	// Vérifier que le fichier existe et est chiffré
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Erreur lecture fichier: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Le fichier est vide")
	}

	// Le contenu ne doit pas être du JSON lisible (doit être chiffré)
	if string(data[:1]) == "{" {
		t.Fatal("Le fichier n'est pas chiffré !")
	}

	// Test 2: Récupérer tous les RSVPs
	rsvps, err := storage.FindAll()
	if err != nil {
		t.Fatalf("Erreur FindAll: %v", err)
	}
	if len(rsvps) != 1 {
		t.Fatalf("Attendu 1 RSVP, obtenu %d", len(rsvps))
	}
	if rsvps[0].ID != "test-id-1" {
		t.Errorf("Mauvais ID: attendu 'test-id-1', obtenu '%s'", rsvps[0].ID)
	}

	// Test 3: Ajouter un deuxième RSVP
	rsvp2, _ := domain.NewRSVP("Marie", "Martin", true, 1, 0, "", "")
	rsvp2.ID = "test-id-2"
	if err := storage.Save(rsvp2); err != nil {
		t.Fatalf("Erreur sauvegarde RSVP 2: %v", err)
	}

	rsvps, err = storage.FindAll()
	if err != nil {
		t.Fatalf("Erreur FindAll: %v", err)
	}
	if len(rsvps) != 2 {
		t.Fatalf("Attendu 2 RSVPs, obtenu %d", len(rsvps))
	}

	// Test 4: Récupérer par ID
	found, err := storage.FindByID("test-id-1")
	if err != nil {
		t.Fatalf("Erreur FindByID: %v", err)
	}
	if found.FirstName != "Jean" {
		t.Errorf("Mauvais prénom: attendu 'Jean', obtenu '%s'", found.FirstName)
	}

	// Test 5: ID inexistant
	_, err = storage.FindByID("inexistant")
	if err != ErrNotFound {
		t.Errorf("Attendu ErrNotFound, obtenu: %v", err)
	}
}

func TestInvalidKey(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	// Test avec une clé invalide (pas base64)
	_, err := NewEncryptedFileStorage(filePath, "invalid-key!!!")
	if err != ErrInvalidKey {
		t.Errorf("Attendu ErrInvalidKey, obtenu: %v", err)
	}

	// Test avec une clé de mauvaise taille
	shortKey := base64.StdEncoding.EncodeToString([]byte("short"))
	_, err = NewEncryptedFileStorage(filePath, shortKey)
	if err != ErrInvalidKeySize {
		t.Errorf("Attendu ErrInvalidKeySize, obtenu: %v", err)
	}
}
