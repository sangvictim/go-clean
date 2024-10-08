package pkg

import "golang.org/x/crypto/bcrypt"

type Encryption interface {
	Bcrypt(payload string) (string, error)
	CompareHashBrypt(payload string, hash string) bool
}

type EncryptionService struct{}

// NewBcryptService returns a new instance of BcryptService.
func NewBcryptService() Encryption {
	return &EncryptionService{}
}

// Bcrypt method hashes the input string using bcrypt and returns the hash or error.
func (s *EncryptionService) Bcrypt(payload string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(payload), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHashBrypt method compares a bcrypt hash with the input string.
func (s *EncryptionService) CompareHashBrypt(payload string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(payload))
	return err == nil
}
