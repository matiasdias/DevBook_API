package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Login é responsalvel por autenticar um usuario na API
func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	var usuario modelos.Usuarios //vindo na requisição
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, erro := repositorio.BuscarPorEmail(ctx, usuario.Email) // usuario salvo no banco
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//esta verificando com a senha que está no banco
	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro) // não está autorizado a fazer essa operação
		return
	}

	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	usuarioID := strconv.FormatUint(usuarioSalvoNoBanco.ID, 10)

	respostas.JSON(w, http.StatusOK, modelos.DadosAutenticacao{ID: usuarioID, Token: token})
}
