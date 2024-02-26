package modelos

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario representa um usuario usando a rede social
type Usuarios struct {
	ID       uint64    `json:"id,omitempty"` // se o id tiver em branco ele não passsa
	Nome     string    `json:"nome,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoem,omitempty"`
}

// Preparar vai chamar os metodos para validar e formatar o usuario recebido
func (Usuarios *Usuarios) Preparar(etapa string) error {
	if erro := Usuarios.validar(etapa); erro != nil {
		return erro
	}
	if erro := Usuarios.formatar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (Usuarios *Usuarios) validar(etapa string) error {
	if Usuarios.Nome == "" {
		return errors.New("O nome é obrigatorio e não pode estar em branco")
	}
	if Usuarios.Nick == "" {
		return errors.New("O Nick é obrigatorio e não pode estar em branco")
	}
	if Usuarios.Email == "" {
		return errors.New("O Email é obrigatorio e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(Usuarios.Email); erro != nil {
		return errors.New("O e-mail inserido é invalido")
	}

	if etapa == "cadastro" && Usuarios.Senha == "" {
		return errors.New("A senha é obrigatorio e não pode estar em branco")
	}
	return nil
}

func (Usuarios *Usuarios) formatar(etapa string) error {
	Usuarios.Nome = strings.TrimSpace(Usuarios.Nome)
	Usuarios.Nick = strings.TrimSpace(Usuarios.Nick)
	Usuarios.Email = strings.TrimSpace(Usuarios.Email)

	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(Usuarios.Senha)
		if erro != nil {
			return erro
		}
		Usuarios.Senha = string(senhaComHash) // convertida para string
	}
	return nil
}
