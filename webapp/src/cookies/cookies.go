package cookies

import (
	"net/http"
	"time"
	"webapp/src/config"

	"github.com/gorilla/securecookie"
)

var s *securecookie.SecureCookie

// Configurar utiliza as variaveis do ambiente para a criação do SecureCookie
func Configurar() {
	s = securecookie.New(config.HashKey, config.BlocKey)
}

// Salvar registra as informaçoes de autenticacao
func Salvar(w http.ResponseWriter, ID, token string) error {
	dados := map[string]string{
		"id":    ID,
		"token": token,
	}
	dadosCodificados, erro := s.Encode("dados", dados) // autenticado e criptografado
	if erro != nil {
		return erro
	}
	// colocar os cokies no browse
	http.SetCookie(w, &http.Cookie{
		Name:     "dados",
		Value:    dadosCodificados,
		Path:     "/",  // funciona em todos os lugares
		HttpOnly: true, // eliminar o risco do cookie ser acessado pelo cliente
	})

	return nil
}

// ler retorna os valores armazenados no cookie
func Ler(r *http.Request) (map[string]string, error) {
	cookie, erro := r.Cookie("dados") //lendo os cookies codificados
	if erro != nil {
		return nil, erro
	}
	valores := make(map[string]string) // alocando um map vazio na momoria e jogando na variavel
	if erro = s.Decode("dados", cookie.Value, &valores); erro != nil {
		return nil, erro
	}
	return valores, nil
}

// Deletar apaga os valortes dentro do cookies da aplicacao ao deslogar
func Deletar(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "dados",
		Value:    "",
		Path:     "/",             // funciona em todos os lugares
		HttpOnly: true,            // eliminar o risco do cookie ser acessado pelo cliente
		Expires:  time.Unix(0, 0), // tempo de expiração do cookies
	})
}
