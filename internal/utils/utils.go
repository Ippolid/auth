package utils

import "golang.org/x/crypto/bcrypt"

// VerifyPassword проверяет, соответствует ли введенный пароль хешированному паролю
func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}
