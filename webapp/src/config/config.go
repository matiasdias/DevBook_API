package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	//APIURL representa a URL para comunicacao coma API
	APIURL = ""
	//Porta onde a aplicacao esta rodando
	Porta = 0
	//HasKey é utilizada para autenticacar o cookie
	HashKey []byte
	//BlocKey é utilizada para criptografar os dados do cookie
	BlocKey []byte
)

// Carregar inializa as variaveis de ambiente
func Carregar() {
	var erro error
	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("APP_PORT"))
	if erro != nil {
		log.Fatal(erro)
	}

	APIURL = os.Getenv("API_URL")
	HashKey = []byte(os.Getenv("HASH_KEY"))
	BlocKey = []byte(os.Getenv("BLOCK_KEY"))

}
