package utils

import "golang.org/x/crypto/bcrypt"

// SenhaHash usa bcrypt para encriptar a senha
func SenhaHash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

// SenhaHashCheck confirma se a senha encriptada é correta
func SenhaHashCheck(SenhaHash, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(SenhaHash), []byte(senha))
}
