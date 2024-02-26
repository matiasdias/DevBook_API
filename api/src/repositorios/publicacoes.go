package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type Publicacoes struct {
	db *sql.DB
}

//NovoRepositorioDeUsuarios cria um repositorio de usuario
func NovoRepositorioDePublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

//Criar uma nova publicacao no BANCO de dados
func (repositorio Publicacoes) Criar(publicacoes modelos.Publicacao) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"insert into publicacoes (titulo, conteudo, autor_id) values (?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()
	resultado, erro := statement.Exec(publicacoes.Titulo, publicacoes.Conteudo, publicacoes.AutorID)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}
	return uint64(ultimoIDInserido), nil
}

//BuscarPorID traz uma unica publicacao no banco de dados
func (repositorio Publicacoes) BuscarPorID(publicacaoID uint64) (modelos.Publicacao, error) {
	linha, erro := repositorio.db.Query(
		`select p.*, u.nick from
		publicacoes p inner join usuarios u
		on u.id = p.autor_id where p.id = ?`,
		publicacaoID,
	)
	if erro != nil {
		return modelos.Publicacao{}, erro
	}
	defer linha.Close()

	var publicacao modelos.Publicacao

	if linha.Next() {
		if erro = linha.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); erro != nil {
			return modelos.Publicacao{}, erro
		}

	}
	return publicacao, nil
}

//Buscar traz as publicacoes dos usuarios seguidos e tambem do usuario que fe\ a propria requisição
func (repositorio Publicacoes) Buscar(usuarioID uint64) ([]modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		select distinct p.*, u.nick from publicacoes p
		inner join usuarios u on u.id = p.autor_id
		inner join seguidores s on p.autor_id = s.usuario_id
		where u.id = ? or s.seguidor_id = ?`,
		usuarioID, usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()
	var publicacoes []modelos.Publicacao

	for linhas.Next() {
		var publicacao modelos.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}
		publicacoes = append(publicacoes, publicacao)
	}
	return publicacoes, nil
}

// Atualizar altera os dados de uma publicacao
func (repositorio Publicacoes) Atualizar(publicacaoID uint64, publicacao modelos.Publicacao) error {
	statement, erro := repositorio.db.Prepare(
		"update publicacoes set titulo = ?, conteudo = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacaoID); erro != nil {
		return erro
	}
	return nil
}

//Deletar exclui uma publicacao do banco de dados
func (repositorio Publicacoes) Deletar(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare("delete from publicacoes where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}
	return nil
}

//BuscarPorUsuario traz todas as publicacoes de um usuario especifico
func (repositorio Publicacoes) BuscarPorUsuario(usuarioID uint64) ([]modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`select p.*, u.nick from publicacoes p
	join usuarios u on u.id = p.autor_id where p.autor_id = ?`, usuarioID)
	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()
	var publicacoes []modelos.Publicacao

	for linhas.Next() {
		var publicacao modelos.Publicacao

		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publicacao)
	}
	return publicacoes, nil
}

//Curtir adicona uma curtida na publicacao
func (repositorio Publicacoes) Curtir(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare("update publicacoes set curtidas = curtidas + 1 where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}
	return nil
}

func (repositorio Publicacoes) Descurtir(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare(`
	update publicacoes set curtidas =
	CASE
		WHEN curtidas > 0 THEN curtidas -1 ELSE 0
	END where id = ? `)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}
	return nil
}
