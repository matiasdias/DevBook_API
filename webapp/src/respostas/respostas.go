package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

//ErroAPI representa a resposta de erro da API
type ErroAPI struct {
	Erro string `json:"erro"`
}

//JSON retorna uma resposta em formato json
func JSON(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode != http.StatusNoContent {
		if erro := json.NewEncoder(w).Encode(dados); erro != nil {
			log.Fatal(erro) // <- so tinha esse if erro
		}
	}
}

//TratarStatusCodeDeErro trata os erros
func TratarStatusCodeDeErro(w http.ResponseWriter, r *http.Response) {
	var erro ErroAPI
	json.NewDecoder(r.Body).Decode(&erro)
	JSON(w, r.StatusCode, erro)
}
