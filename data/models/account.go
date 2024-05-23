package models

import "time"

type Tabler interface {
	TableName() string
}

type Account struct {
	ID                 uint                `gorm:"primaryKey" json:"id"`
	Number             string              `json:"number"`
	DateOpened         time.Time           `json:"date_opened"`
	Balance            float64             `json:"balance"`
	TypeID             uint                `json:"type_id"`
	CreditAccountDatas []CreditAccountData `gorm:"foreignKey:AccountID"`
}

func (Account) TableName() string {
	return "account"
}
