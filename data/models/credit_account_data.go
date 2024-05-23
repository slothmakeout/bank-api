package models

import (
	"time"
)

type CreditAccountData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
	Debtor    bool      `gorm:"not null" json:"debtor"`
	AccountID uint      `gorm:"not null;index;" json:"account_id"`
	Fee       float64   `gorm:"not null" json:"fee"`
	Debt      float64   `gorm:"not null" json:"debt"`
	Account   *Account  `gorm:"foreignKey:AccountID"`
}
