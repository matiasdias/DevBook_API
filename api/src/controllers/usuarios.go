package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuarios cria um novo usuario
func CriarUsuarios(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuarios
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
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
	NewUsuario, err := repositorio.Criar(ctx, usuario)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	responseJSON, err := json.Marshal(NewUsuario)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

// BuscarUsusarios busca um usuario por parametro
func BuscarUsusarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario")) //TOlower -: letras minusculas
	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUsusario busca o usuario pelo id  // OK
func BuscarUsusario(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close() // adia a execução dessa função

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarPorId(ctx, usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if usuarios.ID == uint64(0) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(usuarios)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// AtualizarUsuarios Atualiza os dados do usuario // OK
func AtualizarUsuarios(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	//usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusUnauthorized, erro)
	//	return
	//}
	//if usuarioID != usuarioIDNoToken {
	//	respostas.Erro(w, http.StatusForbidden, errors.New(
	//		"Não é possivel atualizar o usuario que não seja o seu ")) //proibido de fazer isso
	//	return
	//}

	//fmt.Println(usuarioIDNoToken)

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuarios
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edição"); erro != nil {
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
	existingUser, err := repositorio.BuscarPorId(ctx, usuarioID)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Verifica se esse funcionario existe
	if existingUser.ID == uint64(0) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verifica se o email existe ou não no banco de dados
	if exists, err := repositorio.FindByEmailExists(ctx, usuario.Email); err != nil {
		http.Error(w, "Error checking email existence", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Atualiza o usuário
	if erro = repositorio.Atualizar(ctx, usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	message := "User updated successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// DeletarUsuarios Remove um usuario do banco de dados // OK
func DeletarUsuarios(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	//usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusUnauthorized, erro)
	//	return
	//}
	//if usuarioID != usuarioIDNoToken {
	//	respostas.Erro(w, http.StatusForbidden, errors.New(
	//		"Não é possivel deletar o usuario que não seja o seu ")) //proibido de fazer isso
	//	return
	//}

	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorios := repositorios.NovoRepositorioDeUsuarios(db)

	existingUser, err := repositorios.BuscarPorId(ctx, usuarioID)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Verifica se esse funcionario existe
	if existingUser.ID == uint64(0) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if erro = repositorios.Deletar(ctx, usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	message := "Employee deleted successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// SeguirUsuario permite que um usuario siga outro
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusUnauthorized, erro)
	//	return
	//}

	var usuario modelos.Usuarios
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	//if seguidorID == usuarioID {
	//	respostas.Erro(w, http.StatusForbidden, errors.New("Não é possivel seguir voce mesmo"))
	//	return
	//}
	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorios := repositorios.NovoRepositorioDeUsuarios(db)

	seguidorID, erro := repositorios.BuscarPorEmail(ctx, usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if seguidorID.ID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possivel seguir voce mesmo"))
		return
	}
	if erro = repositorios.Seguir(ctx, usuarioID, seguidorID.ID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// NotFollow permite que o usuario pare de seguir outro usuario
func NotFollow(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possivel parar de seguir voce mesmo"))
		return
	}

	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorios := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorios.NotFollow(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarSeguidores traz todos os seguidores de um usuario
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorios := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorios.BuscaSeguidores(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores)
}

// BuscarSeguindo traz todos os usuarios que um determinado usuario está seguindo
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorios := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorios.Seguindo(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarSenha permite alterar a senha de um usuario
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDNoToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New(
			"Não é possivel atualziar a senha de um usuario que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorios := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalvaNoBanco, erro := repositorios.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New(
			"A senha atual não condiz com a que está salva no banco"))
		return
	}
	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = repositorios.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusNoContent, nil)

}
