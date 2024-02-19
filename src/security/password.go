package security

import "golang.org/x/crypto/bcrypt"

func GenerateHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
}
