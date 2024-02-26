package modelos

import "time"

//Publicacao cordo da publicacao
type Publicacao struct {
	ID        uint64    `json:"id,omitemply"`
	Titulo    string    `json:"titulo,omitemply"`
	Conteudo  string    `json:"conteudo,omitemply"`
	AutorID   uint64    `json:"autorId,omitemply"`
	AutorNick string    `json:"autorNick,omitemply"`
	Curtidas  uint64    `json:"curtidas"`
	CriadaEm  time.Time `json:"criadaEm,omitemply"`
}
