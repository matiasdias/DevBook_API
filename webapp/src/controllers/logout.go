package controllers

import (
	"net/http"
	"webapp/src/cookies"
)

//Logout remove os dados da aplicacao salva no browser do usuario
func Logout(w http.ResponseWriter, r *http.Request) {
	cookies.Deletar(w)
	http.Redirect(w, r, "/login", 302)
}
