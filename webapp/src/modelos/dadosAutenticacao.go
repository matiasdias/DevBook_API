package modelos

//DadosAutenticacao contem o id e o token para autenticar o login
type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}
