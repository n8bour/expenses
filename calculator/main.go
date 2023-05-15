package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/n8bour/expenses/calculator/api"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/middleware"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("./calculator/.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	addr := os.Getenv("ADDR")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PW")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	pdb, err := sql.Open("postgres", psqlInfo)
	pdb.SetConnMaxLifetime(time.Minute * 3)
	pdb.SetMaxOpenConns(10)
	pdb.SetMaxIdleConns(10)
	defer func() {
		err := pdb.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewSqlExpenseStore(pdb)
	calculatorService := internal.NewExpenseService(store)
	handleCalculator := api.NewHandleCalculator(calculatorService)

	router := httprouter.New()

	router.POST("/expense", api.WrapHandlers(handleCalculator.HandlePostCalculation))
	router.GET("/expense/:id", api.WrapHandlers(handleCalculator.HandleGetCalculation))

	router.GET("/hello", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		}
		_ = json.NewEncoder(w).Encode("Hello World")
	})

	log.Printf("Server is up and running on: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.SimpleLogging(router)))
}
