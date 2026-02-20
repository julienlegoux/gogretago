package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/lgxju/gogretago/internal/domain/services"
	"golang.org/x/crypto/argon2"
)

// Argon2 configuration matching the TypeScript implementation
const (
	argonTime    = 2
	argonMemory  = 19 * 1024 // 19456 KiB
	argonThreads = 1
	argonKeyLen  = 32
	argonSaltLen = 16
)

// ArgonPasswordService implements PasswordService using Argon2id
type ArgonPasswordService struct{}

// NewArgonPasswordService creates a new ArgonPasswordService
func NewArgonPasswordService() services.PasswordService {
	return &ArgonPasswordService{}
}

// Hash hashes a password using Argon2id
func (s *ArgonPasswordService) Hash(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Encode in the standard format: $argon2id$v=19$m=19456,t=2,p=1$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads, b64Salt, b64Hash)

	return encoded, nil
}

// Verify verifies a password against a hash
func (s *ArgonPasswordService) Verify(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Compute hash of provided password
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expectedHash)))

	// Constant-time comparison
	return subtle.ConstantTimeCompare(expectedHash, computedHash) == 1, nil
}
