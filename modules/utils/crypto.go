package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"

	"golang.org/x/crypto/argon2"
)

type AuthConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen int32
}

func GenerateSalt(length int32) ([]byte, error) {
	salt := make([]byte, length)

	_, err := rand.Read(salt)

	return salt, err
}

func GeneratePasswordHash(password string, salt []byte, config AuthConfig) []byte {
	return argon2.Key([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)
}

func GenerateEncodedSaltAndPasswordHash(password string, config AuthConfig) (encodedPassword string, encodedSalt string, err error) {
	salt, err := GenerateSalt(config.SaltLen)

	if err != nil {
		return encodedPassword, encodedSalt, err
	}

	hashedPassword := GeneratePasswordHash(password, salt, config)

	encodedSalt = base32.StdEncoding.EncodeToString(salt[:])
	encodedPassword = base32.StdEncoding.EncodeToString(hashedPassword[:])

	return encodedPassword, encodedSalt, err
}

func DecodeAndComparePasswords(plainTextPassword string, encodedPassword string, encodedSalt string, config AuthConfig) (bool, error) {
	salt, err := base32.StdEncoding.DecodeString(encodedSalt)

	if err != nil {
		return false, err
	}

	hashedPassword, err := base32.StdEncoding.DecodeString(encodedPassword)

	if err != nil {
		return false, err
	}

	hashedComparison := GeneratePasswordHash(plainTextPassword, salt, config)

	return bytes.Equal(hashedComparison, hashedPassword), nil
}
