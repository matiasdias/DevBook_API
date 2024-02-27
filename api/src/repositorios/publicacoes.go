package repositorios

import (
	"api/src/modelos"
	"context"
	"database/sql"
	"log"
)

type Publicacoes struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositorio de usuario
func NovoRepositorioDePublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

// Criar uma nova publicacao no banco de dados
func (repositorio Publicacoes) Criar(ctx context.Context, publicacoes modelos.Publicacao) (uint64, error) {
	insertQuery := "INSERT INTO public.publicacoes (titulo, conteudo, autor_id) VALUES ($1, $2, $3) RETURNING id"

	row := repositorio.db.QueryRowContext(ctx, insertQuery, publicacoes.Titulo, publicacoes.Conteudo, publicacoes.AutorID)
	var publicacaoID uint64

	err := row.Scan(&publicacaoID)
	if err != nil {
		return 0, err
	}

	return publicacaoID, nil
}

// BuscarPorID traz uma unica publicacao no banco de dados //OK
func (repositorio Publicacoes) BuscarPorID(ctx context.Context, publicacaoID uint64) (modelos.Publicacao, error) {
	query := "SELECT p.id, p.titulo, p.conteudo, p.autor_id, p.curtidas, p.criado_em, u.nick FROM public.publicacoes p INNER JOIN public.usuarios u ON u.id = p.autor_id WHERE p.id = $1"

	var public modelos.Publicacao
	err := repositorio.db.QueryRowContext(ctx, query, publicacaoID).
		Scan(&public.ID,
			&public.Titulo,
			&public.Conteudo,
			&public.AutorID,
			&public.Curtidas,
			&public.CriadaEm,
			&public.AutorNick,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error list publicaçao with ID %d: %v", publicacaoID, err)
			return modelos.Publicacao{}, nil
		}
	}

	return public, nil
}

// Buscar traz as publicacoes dos usuarios seguidos e tambem do usuario que fe\ a propria requisição // OK
func (repositorio Publicacoes) Buscar() ([]modelos.Publicacao, error) {
	query := "SELECT p.*, u.nick from public.publicacoes p INNER JOIN public.usuarios u ON u.id = p.autor_id ORDER BY id ASC"

	rows, err := repositorio.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publicacoes []modelos.Publicacao

	for rows.Next() {
		var publicacao modelos.Publicacao

		if err := rows.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); err != nil {
			return nil, err
		}
		publicacoes = append(publicacoes, publicacao)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return publicacoes, nil

}

// Atualizar altera os dados de uma publicacao //OK
func (repositorio Publicacoes) Atualizar(ctx context.Context, publicacaoID uint64, publicacao modelos.Publicacao) error {
	query := "UPDATE public.publicacoes SET titulo = $1, autor_id = $2, conteudo = $3 WHERE id = $4"

	_, err := repositorio.db.ExecContext(ctx, query,
		publicacao.Titulo,
		publicacao.AutorID,
		publicacao.Conteudo, publicacaoID)
	if err != nil {
		log.Printf("Error updating publicacao with ID %d: %v", publicacaoID, err)
		return err
	}

	return nil
}

// Deletar remove uma publicacao do banco de dados // OK
func (repositorio Publicacoes) Deletar(ctx context.Context, publicacaoID uint64) error {
	query := "DELETE FROM public.publicacoes WHERE id = $1"

	_, err := repositorio.db.ExecContext(ctx, query, publicacaoID)
	if err != nil {
		log.Printf("Error delete publicação with ID %d: %v", publicacaoID, err)
		return err
	}

	return nil
}

// BuscarPorUsuario traz todas as publicacoes de um usuario especifico
func (repositorio Publicacoes) BuscarPorUsuario(ctx context.Context, usuarioID uint64) ([]modelos.Publicacao, error) {
	query := "SELECT p.id, p.titulo, p.conteudo, p.autor_id, p.curtidas, p.criado_em, u.nick FROM public.publicacoes p INNER JOIN public.usuarios u ON u.id = p.autor_id WHERE p.autor_id = $1"

	rows, err := repositorio.db.QueryContext(ctx, query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publicacoes []modelos.Publicacao
	for rows.Next() {
		var publicacao modelos.Publicacao
		if err := rows.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadaEm,
			&publicacao.AutorNick,
		); err != nil {
			return nil, err
		}

		publicacoes = append(publicacoes, publicacao)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return publicacoes, nil
}

// Curtir Adiciona uma curtida a uma publicação
func (repositorio Publicacoes) Curtir(ctx context.Context, publicacaoID uint64) error {

	query := "UPDATE public.publicacoes SET curtidas = curtidas + 1 WHERE id = $1"

	statement, erro := repositorio.db.PrepareContext(ctx, query)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.ExecContext(ctx, publicacaoID); erro != nil {
		return erro
	}
	return nil
}

// Descurtir Responsável por descurti uma publicação // OK
func (repositorio Publicacoes) Descurtir(ctx context.Context, publicacaoID uint64) error {
	query := "UPDATE public.publicacoes SET curtidas = CASE WHEN curtidas > 0 THEN curtidas - 1 ELSE 0 END WHERE id = $1"

	statement, erro := repositorio.db.PrepareContext(ctx, query)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.ExecContext(ctx, publicacaoID); erro != nil {
		return erro
	}
	return nil
}
