package utils

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
)

func HashPassword(uuid uuid.UUID, password string) string {
	hash256 := sha256.New()
	salted := fmt.Sprintf("%s%s", uuid.String(), password)
	hash256.Write([]byte(salted))
	hash := fmt.Sprintf("%x", hash256.Sum(nil))
	return hash
}
