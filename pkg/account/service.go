package accountService

import (
    "sync"

    "bank-api/data/models"
    "gorm.io/gorm"
)

type AccountService struct {
    db *gorm.DB
}

var accountServiceInstance *AccountService
var accountServiceOnce sync.Once

func GetService(db *gorm.DB) *AccountService {
    accountServiceOnce.Do(func() {
        accountServiceInstance = NewAccountService(db)
    })
    return accountServiceInstance
}

func NewAccountService(db *gorm.DB) *AccountService {
    return &AccountService{
        db: db,
    }
}

func (as *AccountService) GetAllAccounts() ([]*models.Account, error) {
    var accounts []*models.Account
    result := as.db.Find(&accounts)
    if result.Error != nil {
        return nil, result.Error
    }
    return accounts, nil
}

func (as *AccountService) GetAccountById(id int) (*models.Account, error) {
    var account models.Account
    result := as.db.First(&account, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &account, nil
}

func (as *AccountService) AddAccount(account *models.Account) error {
	result := as.db.Create(account)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (as *AccountService) UpdateAccount(id int, updatedAccount *models.Account) error {
	// Поиск счета по ID
	existingAccount, err := as.GetAccountById(id)
	if err != nil {
		return err
	}

	// Обновление полей счета
	existingAccount.Number = updatedAccount.Number
	existingAccount.DateOpened = updatedAccount.DateOpened
	existingAccount.Balance = updatedAccount.Balance
	existingAccount.TypeID = updatedAccount.TypeID

	// Сохранение обновленного счета
	result := as.db.Save(existingAccount)
	if result.Error != nil {
		return result.Error
	}
	
	return nil
}

func (as *AccountService) DeleteAccount(id int) error {
	// Поиск счета по ID
	existingAccount, err := as.GetAccountById(id)
	if err != nil {
		return err
	}

	// Удаление счета
	result := as.db.Delete(existingAccount)
	if result.Error != nil {
		return result.Error
	}

	return nil
}



