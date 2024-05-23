package creditAccountDataService

import (
    "sync"

    "bank-api/data/models"
    "gorm.io/gorm"
)

type CreditAccountDataService struct {
    db *gorm.DB
}

var creditAccountDataServiceInstance *CreditAccountDataService
var creditAccountDataServiceOnce sync.Once

func GetService(db *gorm.DB) *CreditAccountDataService {
    creditAccountDataServiceOnce.Do(func() {
        creditAccountDataServiceInstance = NewCreditAccountDataService(db)
    })
    return creditAccountDataServiceInstance
}

func NewCreditAccountDataService(db *gorm.DB) *CreditAccountDataService {
    return &CreditAccountDataService{
        db: db,
    }
}

func (cas *CreditAccountDataService) GetAllCreditAccounts() ([]*models.CreditAccountData, error) {
    var creditAccounts []*models.CreditAccountData
    result := cas.db.Find(&creditAccounts)
    if result.Error != nil {
        return nil, result.Error
    }
    return creditAccounts, nil
}

func (cas *CreditAccountDataService) GetCreditAccountById(id int) (*models.CreditAccountData, error) {
    var creditAccount models.CreditAccountData
    result := cas.db.First(&creditAccount, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &creditAccount, nil
}

func (cas *CreditAccountDataService) AddCreditAccount(creditAccount *models.CreditAccountData) error {
    result := cas.db.Create(creditAccount)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (cas *CreditAccountDataService) UpdateCreditAccount(id int, updatedCreditAccount *models.CreditAccountData) error {
    // Поиск кредитного счета по ID
    existingCreditAccount, err := cas.GetCreditAccountById(id)
    if err != nil {
        return err
    }

    // Обновление полей кредитного счета
    existingCreditAccount.EndDate = updatedCreditAccount.EndDate
    existingCreditAccount.Debtor = updatedCreditAccount.Debtor
    existingCreditAccount.AccountID = updatedCreditAccount.AccountID
    existingCreditAccount.Fee = updatedCreditAccount.Fee
    existingCreditAccount.Debt = updatedCreditAccount.Debt

    // Сохранение обновленного кредитного счета
    result := cas.db.Save(existingCreditAccount)
    if result.Error != nil {
        return result.Error
    }

    return nil
}

func (cas *CreditAccountDataService) DeleteCreditAccount(id int) error {
    // Поиск кредитного счета по ID
    existingCreditAccount, err := cas.GetCreditAccountById(id)
    if err != nil {
        return err
    }

    // Удаление кредитного счета
    result := cas.db.Delete(existingCreditAccount)
    if result.Error != nil {
        return result.Error
    }

    return nil
}

func (cas *CreditAccountDataService) GetAccountWithCreditData(accountID uint) (*models.CreditAccountData, error) {
    var creditAccountData models.CreditAccountData

    result := cas.db.Preload("Account").Where("account_id = ?", accountID).First(&creditAccountData)

    if result.Error != nil {
        return nil, result.Error
    }

    return &creditAccountData, nil
}