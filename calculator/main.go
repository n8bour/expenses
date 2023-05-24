package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/n8bour/expenses/calculator/api"
	"github.com/n8bour/expenses/calculator/db"
	"github.com/n8bour/expenses/calculator/internal"
	"github.com/n8bour/expenses/calculator/middleware"
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

	userStore := db.NewSQLUserStore(pdb)
	expenseStore := db.NewSQLExpenseStore(pdb)

	calculatorService := internal.NewExpenseService(expenseStore)
	userService := internal.NewUserService(userStore, expenseStore)

	handleCalculator := api.NewHandleCalculator(calculatorService)
	handleUser := api.NewHandleUser(userService)
	handleAuth := api.NewAuthHandler(userStore)

	router := chi.NewRouter()
	router.Use(middleware.SimpleLogger)
	router.Post("/login", api.WrapHandlers(handleAuth.HandleAuth))
	router.Post("/user", api.WrapHandlers(handleUser.HandlePostUser))

	apiV1 := router.With(middleware.JwtAuth).Group(func(r chi.Router) {
		r.Post("/expense", api.WrapHandlers(handleCalculator.HandlePostCalculation))
		r.Get("/expense/:id", api.WrapHandlers(handleCalculator.HandleGetCalculation))
		r.Get("/expense", api.WrapHandlers(handleCalculator.HandleListCalculation))
		r.Get("/user/:id", api.WrapHandlers(handleUser.HandleGetUser))
		r.Get("/user", api.WrapHandlers(handleUser.HandleListUsers))
	})

	router.Mount("/api", apiV1)

	log.Printf("Server is up and running on: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
