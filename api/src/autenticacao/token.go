package autenticacao

import (
	"api/src/banco"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CriarToken retorna um token assinado com as permissões do usuario
func CriarToken(usuairoID uint64) (string, error) {
	permissoes := jwt.MapClaims{}                                    //onde vai conter as permissoes dentro do token
	permissoes["authorized"] = true                                  //usuario está autorizado a utilziar
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()         //duração do token
	permissoes["usuarioId"] = usuairoID                              //usuario que vai utilizar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)   //criar um token
	return token.SignedString([]byte(banco.APIConfigInfo.APISecret)) // assinatura do token
}

// ValidarToken verifica se o token passado é valido
func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return erro
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("Token invalido")
}

// ExtrairUsuarioID retorna o usuarioId que esta salvono token
func ExtrairUsuarioID(r *http.Request) (uint64, error) {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return 0, erro
	}
	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissoes["usuarioId"]), 10, 64)
		if erro != nil {
			return 0, erro
		}
		return usuarioID, nil
	}
	return 0, errors.New("Token invalido")
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Metodo de assinatura inesperado: %v", token.Header)
	}
	return banco.APIConfigInfo.APISecret, nil
}
