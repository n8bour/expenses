package main

import (
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/n8bour/expenses/calculator/api"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load("./calculator/.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	addr := os.Getenv("ADDR")

	pdb, err := db.NewPsqlDB()
	if err != nil {
		log.Fatal(err)
	}

	pdb.SetConnMaxLifetime(time.Minute * 3)
	pdb.SetMaxOpenConns(10)
	pdb.SetMaxIdleConns(10)

	defer func() {
		err := pdb.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	expenseStore := db.NewSqlExpenseStore(pdb)
	calculatorService := internal.NewExpenseService(expenseStore)
	userService := internal.NewUserService(db.NewSqlUserStore(pdb), expenseStore)
	handleCalculator := api.NewHandleCalculator(calculatorService)
	handleUser := api.NewHandleUser(userService)

	router := httprouter.New()

	router.POST("/expense", api.WrapHandlers(handleCalculator.HandlePostCalculation))
	router.GET("/expense/:id", api.WrapHandlers(handleCalculator.HandleGetCalculation))
	router.GET("/expense", api.WrapHandlers(handleCalculator.HandleListCalculation))

	router.POST("/user", api.WrapHandlers(handleUser.HandlePostUser))
	log.Printf("Server is up and running on: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.SimpleLogging(router)))
}
