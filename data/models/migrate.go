package models

import "gorm.io/gorm"

// AutoMigrate выполняет автоматическую миграцию таблиц базы данных.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Account{}, &CreditAccountData{})
}
