package middlewares

import (
	"api/src/autenticacao"
	"log"
	"net/http"
)

// Logger escreve informações da requisicao no terminal
func Logger(ProximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %q %q %q", r.Method, r.RequestURI, r.Host)
		ProximaFuncao(w, r)
	}
}

// Autenticar verifica se o usuario fazendo a requisição está autenticado
func Autenticar(ProximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := autenticacao.ValidarToken(r); erro != nil {
			http.Error(w, "token contains as invalid number", http.StatusUnauthorized)
			return
		}
		ProximaFuncao(w, r)
	}
}
