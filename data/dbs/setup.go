package dbs

import (
	"bank-api/data/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB

func GetDB() *gorm.DB {
	return dbInstance
}

func setupPostgres() (*gorm.DB, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	fmt.Println(dbHost)
	connectionString := fmt.Sprintf("%s://%s:%s@localhost:%s/%s",
		dbHost,
		dbUser,
		dbPassword,
		dbPort,
		dbName)
	fmt.Println(connectionString)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitializeDatabaseLayer() error {

	var db *gorm.DB
	var err error

	db, err = setupPostgres()

	if err != nil {
		return err
	}

	err = models.AutoMigrate(db)
	if err != nil {
		return err
	}
	dbInstance = db
	fmt.Println("migration completed")
	return nil
}
