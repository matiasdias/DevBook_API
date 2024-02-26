package seguranca

import "golang.org/x/crypto/bcrypt"

//Hash recebe uma string e coloca um hash nela
func Hash(senha string) ([]byte, error) { // Gerar um hash a partir de uma senha
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

//VerificarSenha compara uma senha com o hash e retorna se elas são iguais
func VerificarSenha(senhaComHash, senhaString string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaComHash), []byte(senhaString))
}
