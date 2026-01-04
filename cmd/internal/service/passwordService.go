package service

import "golang.org/x/crypto/bcrypt"

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),     // senha salva no banco (hash)
		[]byte(password), // senha digitada pelo usuário
	)
	return err == nil
}
