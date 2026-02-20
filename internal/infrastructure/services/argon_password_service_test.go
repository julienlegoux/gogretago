package services

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash_ProducesArgon2idFormat(t *testing.T) {
	svc := NewArgonPasswordService()
	hash, err := svc.Hash("TestPassword123")
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(hash, "$argon2id$v="), "hash should start with $argon2id$v=")
}

func TestHash_DifferentSalts(t *testing.T) {
	svc := NewArgonPasswordService()
	hash1, err := svc.Hash("SamePassword1")
	require.NoError(t, err)

	hash2, err := svc.Hash("SamePassword1")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "two hashes of the same password should differ due to random salt")
}

func TestVerify_CorrectPassword(t *testing.T) {
	svc := NewArgonPasswordService()
	password := "CorrectPassword1"

	hash, err := svc.Hash(password)
	require.NoError(t, err)

	ok, err := svc.Verify(password, hash)
	require.NoError(t, err)
	assert.True(t, ok, "correct password should verify successfully")
}

func TestVerify_WrongPassword(t *testing.T) {
	svc := NewArgonPasswordService()

	hash, err := svc.Hash("CorrectPassword1")
	require.NoError(t, err)

	ok, err := svc.Verify("WrongPassword1", hash)
	require.NoError(t, err)
	assert.False(t, ok, "wrong password should not verify")
}

func TestVerify_InvalidHashFormat(t *testing.T) {
	svc := NewArgonPasswordService()
	_, err := svc.Verify("anypassword", "notahash")
	assert.Error(t, err, "invalid hash format should return error")
}

func TestVerify_CorruptedSalt(t *testing.T) {
	svc := NewArgonPasswordService()
	// Valid format but with corrupted base64 salt (invalid characters for base64)
	corrupted := "$argon2id$v=19$m=19456,t=2,p=1$!!!invalid-base64!!!$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	_, err := svc.Verify("anypassword", corrupted)
	assert.Error(t, err, "corrupted salt should return error")
}

func TestVerify_CorruptedHash(t *testing.T) {
	svc := NewArgonPasswordService()
	// Valid format with valid base64 salt but corrupted base64 hash
	corrupted := "$argon2id$v=19$m=19456,t=2,p=1$AAAAAAAAAAAAAAAAAAAAAA$!!!invalid-base64!!!"
	_, err := svc.Verify("anypassword", corrupted)
	assert.Error(t, err, "corrupted hash should return error")
}
