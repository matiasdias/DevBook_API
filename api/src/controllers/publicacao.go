package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"strconv"

	//"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	//"strconv"
	//"github.com/gorilla/mux"
)

// CriarPublicacao Cria uma nova publicacao // OK
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	//if erro != nil {
	//	fmt.Println(usuarioID)
	//	respostas.Erro(w, http.StatusUnauthorized, erro)
	//	return
	//}
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	//publicacao.AutorID = usuarioID

	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao.ID, erro = repositorio.Criar(ctx, publicacao)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//w.Write(publicacao.ID)
}

// BuscarPublicacoes busca todas as publicações no banco de dados // OK
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {
	//usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusUnauthorized, erro)
	//	return
	//}
	db, erro := banco.Connection()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.Buscar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, publicacoes)
}

// BuscarPublicacao busca uma publicacao por id //OK
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao, erro := repositorio.BuscarPorID(ctx, publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, publicacao)
}

// AtualizarPublicacao atualiza as publicacções de acordo com o id forncecido // OK
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusUnauthorized, erro)
	//	return
	//}

	//usuario que está no corpo da requisição
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	//publicacaoSalvaNoBanco, erro := repositorio.BuscarPorID(ctx, publicacaoID)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusInternalServerError, erro)
	//	return
	//}
	//
	//if publicacaoSalvaNoBanco.AutorID != usuarioID {
	//	respostas.Erro(w, http.StatusForbidden, errors.New("Não é possivel atualizar uma publicacao que não seja a sua "))
	//	return
	//}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.Atualizar(ctx, publicacaoID, publicacao); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	message := "publications updated successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// DeletarPublicacao Remove uma publicação // OK
func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusUnauthorized, erro)
	//	return
	//} //usuario que está no corpo da requisição
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	//publicacaoSalvaNoBanco, erro := repositorio.BuscarPorID(ctx, publicacaoID)
	//if erro != nil {
	//	respostas.Erro(w, http.StatusInternalServerError, erro)
	//	return
	//}
	//if publicacaoSalvaNoBanco.AutorID != usuarioID {
	//	respostas.Erro(w, http.StatusForbidden, errors.New("Não é possivel deletar uma publicacao que não seja a sua "))
	//	return
	//}
	if erro = repositorio.Deletar(ctx, publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	message := "publications deleted successfully"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// BuscarPublicacoesPorUsuario traz todas as publicacoes de um usuario especifico //OK
func BuscarPublicacoesPorUsuario(w http.ResponseWriter, r *http.Request) {
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
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.BuscarPorUsuario(ctx, usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.JSON(w, http.StatusOK, publicacoes)
}

// CurtirPublicacao adiciona uma curtida na publicacao // OK
func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	if erro = repositorio.Curtir(ctx, publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	message := "post successfully liked"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// DescurtirPublicacao retira uma curtida na publicacao // OK
func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
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

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	if erro = repositorio.Descurtir(ctx, publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	message := "post successfully disliked"
	responseJSON, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
