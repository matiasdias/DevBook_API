package repositorios

import (
	"api/src/modelos"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositorio de usuario
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

// Criar insere um usuario no banco de dados
func (repositorio Usuarios) Criar(ctx context.Context, usuario modelos.Usuarios) (*modelos.Usuarios, error) {

	insertQuery := "INSERT INTO public.usuarios (nome, nick, email, senha) VALUES ($1, $2, $3, $4) RETURNING id"

	_, err := repositorio.db.ExecContext(ctx, insertQuery, usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if err != nil {
		return nil, err
	}

	var userID uint64
	queryID := "SELECT id FROM public.usuarios WHERE email = $1"
	err = repositorio.db.QueryRowContext(ctx, queryID, usuario.Email).Scan(&userID)
	if err != nil {
		return nil, err
	}

	userResponse := &modelos.Usuarios{
		ID:    userID,
		Nome:  usuario.Nome,
		Nick:  usuario.Nick,
		Email: usuario.Email,
	}

	return userResponse, nil
}

// Buscar traz todos os usuarios que atedem um filtro de nome ou nick
func (repositorio Usuarios) Buscar(nomeOuNick string) ([]modelos.Usuarios, error) {

	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) //%nomeOuNick%

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoem from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOuNick, nomeOuNick,
	)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var usuarios []modelos.Usuarios

	for linhas.Next() {
		var usuario modelos.Usuarios
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarPorId lista um usuário por id // OK
func (repositorio Usuarios) BuscarPorId(ctx context.Context, ID uint64) (modelos.Usuarios, error) {
	query := "SELECT id, nome, nick, email, criadoem FROM public.usuarios WHERE id = $1"

	var usuario modelos.Usuarios
	err := repositorio.db.QueryRowContext(ctx, query, ID).
		Scan(&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error list user with ID %d: %v", ID, err)
			return modelos.Usuarios{}, nil
		}
	}

	return usuario, nil
}

// Atualizar altera as informações de um usuairo no banco de dados //OK
func (repositorio Usuarios) Atualizar(ctx context.Context, ID uint64, usuario modelos.Usuarios) error {
	query := "UPDATE public.usuarios SET nome = $1, nick = $2, email = $3, updateem = current_timestamp WHERE id = $4"

	_, err := repositorio.db.ExecContext(ctx, query,
		usuario.Nome,
		usuario.Nick,
		usuario.Email, ID)
	if err != nil {
		log.Printf("Error updating user with ID %d: %v", ID, err)
		return err
	}

	return nil
}

// Deletar exclui todas as informações de um usuario no banco de dados // OK
func (repositorio Usuarios) Deletar(ctx context.Context, ID uint64) error {
	query := "DELETE FROM public.usuarios WHERE id = $1"

	_, err := repositorio.db.ExecContext(ctx, query, ID)
	if err != nil {
		log.Printf("Error delete user with ID %d: %v", ID, err)
		return err
	}

	return nil
}

// BuscarPorEmail busca um usuario por email e retorna o seu id e senha com hash //OK
func (repositorio Usuarios) BuscarPorEmail(ctx context.Context, email string) (modelos.Usuarios, error) {
	query := "SELECT id, email, senha FROM public.usuarios WHERE email = $1"

	var user modelos.Usuarios
	err := repositorio.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Senha)
	if err != nil {
		return modelos.Usuarios{}, err
	}

	return user, nil
}

// Seguir permite que um usuario siga outro usuário
func (repositorio Usuarios) Seguir(ctx context.Context, usuarioID, seguidorID uint64) error {
	stmt, err := repositorio.db.
		PrepareContext(ctx, `INSERT INTO seguidores (usuario_id, seguidor_id)VALUES ($1, $2)ON CONFLICT DO NOTHING
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, usuarioID, seguidorID)
	if err != nil {
		return err
	}
	return nil
}

// NotFollw permite que o usuario pare de seguir outro
func (repositorio Usuarios) NotFollow(usuarioID, seguidorID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"delete from seguidores where usuario_id = ? and seguidor_id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}
	return nil
}

// BuscaSeguidores traz todos os seguidores de um usuario
func (repositorio Usuarios) BuscaSeguidores(usuarioID uint64) ([]modelos.Usuarios, error) {
	linhas, erro := repositorio.db.Query(
		`select u.id, u.nome, u.nick, u.email, u.criadoEm
		 from usuarios u inner join seguidores s on u.id= s.seguidor_id where s.usuario_id = ?
		 `, usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuarios
	for linhas.Next() {
		var usuario modelos.Usuarios
		if erro = linhas.Scan(
			&usuarioID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

// Seguindo traz todos os usuarios que um determinado usuario está seguindo
func (repositorio Usuarios) Seguindo(usuarioID uint64) ([]modelos.Usuarios, error) {
	linhas, erro := repositorio.db.Query(
		`select u.id, u.nome, u.nick, u.email, u.criadoEm 
		from usuarios u inner join seguidores s on u.id = s.usuario_id where s.seguidor_id = ?
		`, usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuarios
	for linhas.Next() {
		var usuario modelos.Usuarios
		if erro = linhas.Scan(
			&usuarioID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

// BuscarSenha traz a senha de um usuario pelo ID
func (repositorio Usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, erro := repositorio.db.Query("select senha from usuarios where id = ?", usuarioID)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	var usuario modelos.Usuarios
	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "", erro
		}
	}
	return usuario.Senha, nil
}

// AtualizarSenha altera a senha de um usuario
func (repositorio Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("update usuarios set senha = ? where id + ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(senha, usuarioID); erro != nil {
		return erro
	}
	return nil
}

// FindByEmailExists Verifica se o email que tá salvo no banco existe // ok
func (repositorio Usuarios) FindByEmailExists(ctx context.Context, Email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM public.usuarios WHERE email = $1)"

	var exists bool
	err := repositorio.db.QueryRowContext(ctx, query, Email).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if email exists: %v", err)
		return false, err
	}

	return exists, nil
}
