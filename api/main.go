package main

import (
	"api/src/banco"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	banco.Connection()
	r := router.Gerar()

	fmt.Printf("Executando na porta %d\n", banco.APIConfigInfo.APIPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", banco.APIConfigInfo.APIPort), r))
}
