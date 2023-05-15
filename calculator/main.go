package main

import (
	"github.com/joho/godotenv"
	"github.com/n8bour/expenses/calculator/api"
	"log"
	"net/http"
	"os"
)

type Expense struct {
	Type  string
	value float32
}

func main() {
	err := godotenv.Load("./calculator/.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	p := os.Getenv("PORT")

	http.HandleFunc("/calculator", api.NewHandleCalculator())

	log.Printf("Server is up and running on port: %s\n", p)
	log.Fatal(http.ListenAndServe(p, nil))
}
