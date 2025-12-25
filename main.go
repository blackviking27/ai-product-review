package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blackviking27/ai-product-reviwer/internal/api"
	"github.com/joho/godotenv"
)

func main() {

	// Loading env variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	fmt.Print("Started listening on :" + PORT + "\n")
	if err := http.ListenAndServe(":"+PORT, api.NewRouter()); err != nil {
		log.Fatal(err)
	}

}
