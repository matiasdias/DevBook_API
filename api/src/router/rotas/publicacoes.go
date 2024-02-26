package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasPublicacoes = []Rota{
	//crud publicacoes
	{
		URI:                "/publicacoes",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarPublicacao,
		RequerAltenticacao: false,
	},

	{
		URI:                "/publicacoes",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPublicacoes,
		RequerAltenticacao: false,
	},

	{
		URI:                "/publicacoes/{publicacaoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPublicacao,
		RequerAltenticacao: false,
	},

	{
		URI:                "/publicacoes/{publicacaoId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarPublicacao,
		RequerAltenticacao: false,
	},

	{
		URI:                "/publicacoes/{publicacaoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarPublicacao,
		RequerAltenticacao: false,
	},
	//------------------------
	{
		URI:                "/usuarios/{usuarioId}/publicacoes",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPublicacoesPorUsuario,
		RequerAltenticacao: false,
	},

	{
		URI:                "/publicacoes/{publicacaoId}/curtir",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CurtirPublicacao,
		RequerAltenticacao: false,
	},

	{
		URI:                "/publicacoes/{publicacaoId}/descurtir",
		Metodo:             http.MethodPost,
		Funcao:             controllers.DescurtirPublicacao,
		RequerAltenticacao: false,
	},
}
