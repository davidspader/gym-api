package main

import (
	"fmt"
	"gym-api/src/config"
	"gym-api/src/routes"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()
	r := routes.GenerateRouter()

	fmt.Printf("listening on port %d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
