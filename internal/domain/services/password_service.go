package services

// PasswordService defines the interface for password hashing operations
type PasswordService interface {
	Hash(password string) (string, error)
	Verify(password string, hash string) (bool, error)
}
