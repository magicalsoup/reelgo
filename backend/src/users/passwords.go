package users

import (
	"encoding/base64"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

func argon2_hash(password string, salt string) string {
	hashed_pw := argon2.Key([]byte(password), []byte(salt), 3, 32*1024, 4, 32)
	encoded := base64.StdEncoding.EncodeToString(hashed_pw)
	return encoded
}

func getHashedPassword(password string, salt string) string {
	return argon2_hash(password, salt)
}

func generateHashedPassword(password string) (salt string, hashed_password string) {
	generated_salt := uuid.New().String()
	hashed_password = argon2_hash(password, generated_salt)
	return generated_salt, hashed_password
}