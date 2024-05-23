package main

import (
	"bank-api/data/dbs"
	"bank-api/data/models"
	"bank-api/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	accountService "bank-api/pkg/account"
	creditAccountDataService "bank-api/pkg/credit_account_data"

	gohandlers "github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func main() {
	// Инициализация базы данных
	err := dbs.InitializeDatabaseLayer()
	if err != nil {
		panic("Failed to initialize database")
	}

	l := log.New(os.Stdout, "movies-project-microservice-api ", log.LstdFlags)
	db := dbs.GetDB()

	// Инициализация сервисов
	accountService := accountService.NewAccountService(db)
	creditAccountDataService := creditAccountDataService.NewCreditAccountDataService(db)
	// Инициализация обработчиков
	accountHandler := handlers.NewAccountHandler(l, accountService)
	creditAccountDataHandler := handlers.NewCreditAccountDataHandler(l, creditAccountDataService)

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/accounts", accountHandler.GetAllAccounts)
	getRouter.HandleFunc("/accounts/{id:[0-9]+}", accountHandler.GetAccountById)
	getRouter.HandleFunc("/credit_accounts", creditAccountDataHandler.GetAllCreditAccounts)
	getRouter.HandleFunc("/credit_accounts/{id:[0-9]+}", creditAccountDataHandler.GetCreditAccountByID)
	getRouter.HandleFunc("/credit_accounts/full/{id:[0-9]+}", creditAccountDataHandler.GetAccountWithCreditData)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/accounts", accountHandler.AddAccount)
	postRouter.HandleFunc("/credit_accounts", creditAccountDataHandler.AddCreditAccount)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/accounts/{id:[0-9]+}", accountHandler.UpdateAccount)
	putRouter.HandleFunc("/credit_accounts/{id:[0-9]+}", creditAccountDataHandler.UpdateCreditAccount)

	deleteRouter := sm.Methods("DELETE").Subrouter()
	deleteRouter.HandleFunc("/accounts/{id:[0-9]+}", accountHandler.DeleteAccount)
	deleteRouter.HandleFunc("/credit_accounts/{id:[0-9]+}", creditAccountDataHandler.DeleteCreditAccount)

	// CORS
	corsHandler := gohandlers.CORS(
		gohandlers.AllowedOrigins([]string{"*"}),
		gohandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gohandlers.AllowedHeaders([]string{"Content-Type"}),
	)

	s := &http.Server{
		Addr:         "localhost:9090",
		Handler:      corsHandler(sm), // set the default handler
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Создаём канал чтобы принимать сигналы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Запускает сервер в отдельной горутине
	go func() {
		l.Printf("Server listening on %s\n", s.Addr)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("listen: %s\n", err)
		}
	}()

	// Ждём когда сигнал завершит работу сервера
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	// Создаём контекст с таймаутом
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// завершаем работу сервера
	s.Shutdown(tc)
}

func testORM(db *gorm.DB) {
	newAccount := models.Account{
		Number:     "101-221",
		DateOpened: time.Now(),
		Balance:    2000.00,
		TypeID:     1,
	}

	result := db.Create(&newAccount)
	if result.Error != nil {
		panic("Failed to create account")
	}

	fmt.Println("New account created successfully")

	var retrievedAccount models.Account
	db.First(&retrievedAccount, newAccount.ID)

	fmt.Printf("Retrieved account: %+v\n", retrievedAccount)
}
