package modelos

// DadosAutenticacao contem o token e o id do usuario autenticado
type DadosAutenticacao struct {
	ID    uint64 `json:"id"`
	Token string `json:"token"`
}
