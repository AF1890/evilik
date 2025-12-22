package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"wedding-web/internal/domain"
)

var (
	ErrInvalidKey     = errors.New("clé de chiffrement invalide")
	ErrDecryptFailed  = errors.New("échec du déchiffrement")
	ErrNotFound       = errors.New("RSVP non trouvé")
	ErrInvalidKeySize = errors.New("la clé doit faire 32 bytes")
)

// EncryptedFileStorage implémente le stockage chiffré en fichier JSON
type EncryptedFileStorage struct {
	filePath string
	key      []byte
	mu       sync.RWMutex
}

type storageData struct {
	RSVPs []*domain.RSVP `json:"rsvps"`
}

// NewEncryptedFileStorage crée un nouveau storage avec chiffrement AES-GCM
func NewEncryptedFileStorage(filePath string, encryptionKey string) (*EncryptedFileStorage, error) {
	// Décoder la clé base64
	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, ErrInvalidKey
	}

	// Vérifier la taille de la clé (32 bytes pour AES-256)
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	// Créer le répertoire si nécessaire
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	return &EncryptedFileStorage{
		filePath: filePath,
		key:      key,
	}, nil
}

// Save enregistre un RSVP
func (s *EncryptedFileStorage) Save(rsvp *domain.RSVP) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Charger les données existantes
	data, err := s.loadData()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Ajouter le nouveau RSVP
	data.RSVPs = append(data.RSVPs, rsvp)

	// Sauvegarder
	return s.saveData(data)
}

// FindAll retourne tous les RSVPs
func (s *EncryptedFileStorage) FindAll() ([]*domain.RSVP, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := s.loadData()
	if err != nil {
		if os.IsNotExist(err) {
			return []*domain.RSVP{}, nil
		}
		return nil, err
	}

	return data.RSVPs, nil
}

// FindByID retourne un RSVP par son ID
func (s *EncryptedFileStorage) FindByID(id string) (*domain.RSVP, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := s.loadData()
	if err != nil {
		return nil, err
	}

	for _, rsvp := range data.RSVPs {
		if rsvp.ID == id {
			return rsvp, nil
		}
	}

	return nil, ErrNotFound
}

// Delete supprime un RSVP par son ID
func (s *EncryptedFileStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.loadData()
	if err != nil {
		return err
	}

	// Trouver l'index du RSVP à supprimer
	found := false
	newRSVPs := make([]*domain.RSVP, 0, len(data.RSVPs))
	for _, rsvp := range data.RSVPs {
		if rsvp.ID == id {
			found = true
			continue // Ne pas ajouter ce RSVP à la nouvelle liste
		}
		newRSVPs = append(newRSVPs, rsvp)
	}

	if !found {
		return ErrNotFound
	}

	// Sauvegarder les données sans le RSVP supprimé
	data.RSVPs = newRSVPs
	return s.saveData(data)
}

// loadData charge et déchiffre les données
func (s *EncryptedFileStorage) loadData() (*storageData, error) {
	// Lire le fichier
	ciphertext, err := os.ReadFile(s.filePath)
	if err != nil {
		return &storageData{RSVPs: []*domain.RSVP{}}, err
	}

	// Si le fichier est vide, retourner une structure vide
	if len(ciphertext) == 0 {
		return &storageData{RSVPs: []*domain.RSVP{}}, nil
	}

	// Déchiffrer
	plaintext, err := s.decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	// Désérialiser
	var data storageData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// saveData chiffre et sauvegarde les données
func (s *EncryptedFileStorage) saveData(data *storageData) error {
	// Sérialiser
	plaintext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Chiffrer
	ciphertext, err := s.encrypt(plaintext)
	if err != nil {
		return err
	}

	// Écrire dans un fichier temporaire puis renommer (atomic write)
	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, ciphertext, 0600); err != nil {
		return err
	}

	return os.Rename(tmpFile, s.filePath)
}

// encrypt chiffre les données avec AES-GCM
func (s *EncryptedFileStorage) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Générer un nonce aléatoire
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Chiffrer (le nonce est préfixé au ciphertext)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt déchiffre les données avec AES-GCM
func (s *EncryptedFileStorage) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrDecryptFailed
	}

	// Extraire le nonce et le ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Déchiffrer
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}
