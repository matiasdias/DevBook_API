package modelos

import (
	"errors"
	"strings"
	"time"
)

//Publicacao representa uma publicacao de um usuario
type Publicacao struct {
	ID        uint64    `json:"id,omitemply"`
	Titulo    string    `json:"titulo,omitemply"`
	Conteudo  string    `json:"conteudo,omitemply"`
	AutorID   uint64    `json:"autorId,omitemply"`
	AutorNick string    `json:"autorNick,omitemply"`
	Curtidas  uint64    `json:"curtidas"`
	CriadaEm  time.Time `json:"criadaEm,omitemply"`
}

//Preparar vai chamar os metodos para validar e formatar a publicacao recebida
func (publicacao *Publicacao) Preparar() error {
	if erro := publicacao.validar(); erro != nil {
		return erro

	}
	publicacao.formatar()
	return nil
}

//validar verifica se o titulo e o conteudo tem alguma coisa dentro
func (publicacao *Publicacao) validar() error {
	if publicacao.Titulo == "" {
		return errors.New("O titulo é obrigatorio e não pode estar em branco")
	}
	if publicacao.Conteudo == "" {
		return errors.New("O conteudo é obrigatorio e não pode estar em branco")
	}
	return nil
}

// retira todos os espaços em branco da requisição
func (publicacao *Publicacao) formatar() {
	publicacao.Titulo = strings.TrimSpace(publicacao.Titulo)
	publicacao.Conteudo = strings.TrimSpace(publicacao.Conteudo)
}
