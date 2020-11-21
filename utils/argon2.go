package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type Argon2IDConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	Format  string
}

// Returns the default Argon2ID config used for hashing the password.
// For more information on why these parameters are used, see the comments of the argon2.IDKey function. At the time of
// writing this, it was recommended by this draft RFC: https://tools.ietf.org/html/draft-irtf-cfrg-argon2-03#section-9.3
func GetDefaultArgon2IDConfig() *Argon2IDConfig {
	return &Argon2IDConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 4, // TODO: Make this configurable
		KeyLen:  32,
		Format:  "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
	}
}

// Hashes a string using argon2id with the supplied salt and config.
func argon2idHashStringWithSalt(s string, salt []byte, config *Argon2IDConfig) (string, error) {
	hash := argon2idHash([]byte(s), salt, config)

	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	full := fmt.Sprintf(config.Format, argon2.Version, config.Memory, config.Time, config.Threads, b64Salt, b64Hash)
	return full, nil
}

// Uses Argon2ID to hash the supplied string with the supplied config. It generates a random salt. Acts as a public-facing
// wrapper around argon2idHashStringWithSalt.
func Argon2IDHashString(s string, config *Argon2IDConfig) (string, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}
	return argon2idHashStringWithSalt(s, salt, config)
}

// Compares the supplied string with that hash.
func Argon2IDHashCompare(s, hash string) (bool, error) {
	parts := strings.Split(hash, "$")

	// Parse the used parameters
	config := &Argon2IDConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Time, &config.Threads)
	if err != nil {
		return false, nil
	}

	// Parse the salt and the hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	config.KeyLen = uint32(len(decodedHash))

	// Create a comparison hash with the same parameters which should return the same hash
	comparisonHash := argon2.IDKey([]byte(s), salt, config.Time, config.Memory, config.Threads, config.KeyLen)
	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}

// Hashes the supplied input and salt bytes and config with Argon2ID.
func argon2idHash(input, salt []byte, config *Argon2IDConfig) (hash []byte) {
	return argon2.IDKey(input, salt, config.Time, config.Memory, config.Threads, config.KeyLen)
}

// Generates a salt of supplied length.
func generateSalt(n uint) ([]byte, error) {
	salt := make([]byte, n)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}
