package service

import (
	"microservice/db"
	"microservice/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestDeposit(t *testing.T) {
	db, err := db.Connection()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	c := &gin.Context{}
	transactions := []models.Transaction{
		{Type: "deposit",
			Amount:   decimal.NewFromFloat(100),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
		{Type: "deposit",
			Amount:   decimal.NewFromFloat(-1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
		{Type: "z",
			Amount:   decimal.NewFromFloat(1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
	}
	for _, transaction := range transactions {
		t.Run(transaction.Type, func(t *testing.T) {
			wallet := &models.Wallet{}
			err = db.QueryRow("SELECT balance FROM wallet WHERE id_wallet=$1", transaction.WalletId).Scan(&wallet.Balance)
			if err != nil {
				t.Errorf("Error finding wallet: %v", err)
			}
			err = Deposit(db, c, &transaction)
			if transaction.Amount.IsNegative() {
				if err.Error() != "the amount of money to deposit is negative, cannot proceed with transaction" {
					t.Errorf("Expected error to be the amount of money to deposit is negative, cannot proceed with transaction, got: %v", err)
				}
			}
			if transaction.Type != "deposit" {
				if err.Error() != "invalid transaction type" {
					t.Errorf("Expected error to be invalid transaction type but got %v", err)
				}
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	db, err := db.Connection()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	c := &gin.Context{}
	transactions := []models.Transaction{
		{Type: "withdraw",
			Amount:   decimal.NewFromFloat(1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
		{Type: "withdraw",
			Amount:   decimal.NewFromFloat(-1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
		{Type: "z",
			Amount:   decimal.NewFromFloat(1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
	}
	for _, transaction := range transactions {
		t.Run(transaction.Type, func(t *testing.T) {
			wallet := &models.Wallet{}
			err = db.QueryRow("SELECT balance FROM wallet WHERE id_wallet=$1", transaction.WalletId).Scan(&wallet.Balance)
			if err != nil {
				t.Errorf("Error finding wallet: %v", err)
			}
			if transaction.Amount.GreaterThan(wallet.Balance) {
				t.Error("Insufficient balance ")
			}
			err = Withdraw(db, c, &transaction)
			if transaction.Amount.IsNegative() {
				if err.Error() != "the amount of money to withdraw is negative, cannot proceed with transaction" {
					t.Errorf("Expected error to be the amount of money to withdraw is negative, cannot proceed with transaction, got: %v", err)
				}
			}
			if transaction.Type != "withdraw" {
				if err.Error() != "invalid transaction type" {
					t.Errorf("Expected error to be invalid transaction type but got %v", err)
				}
			}
		})
	}
}

func TestMakeTransaction(t *testing.T) {
	db, err := db.Connection()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	c := &gin.Context{}
	transactions := []models.Transaction{
		{Type: "withdraw",
			Amount:   decimal.NewFromFloat(1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
		{Type: "deposit",
			Amount:   decimal.NewFromFloat(1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
		{Type: "z",
			Amount:   decimal.NewFromFloat(1000),
			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
	}
	for _, transaction := range transactions {
		t.Run(transaction.Type, func(t *testing.T) {
			if transaction.Amount.IsNegative() {
				t.Errorf("the amount of money to withdraw is negative, cannot proceed with transaction")
			}
			wallet := &models.Wallet{}
			err = db.QueryRow("SELECT balance FROM wallet WHERE id_wallet=$1", transaction.WalletId).Scan(&wallet.Balance)
			if err != nil {
				t.Errorf("Error finding wallet: %v", err)
			}
			switch transaction.Type {
			case "deposit":
				{
					err = MakeTransaction(db, c, &transaction)
					if err != nil {
						t.Errorf("Expected no error, got: %v", err)
					}
					w := &models.Wallet{}
					db.QueryRow("SELECT * FROM wallet where id_wallet=$1", transaction.WalletId).Scan(&w)
					if w.Balance.Equal(wallet.Balance) {
						t.Errorf("Expected wallet balance to be higher than %v ", wallet.Balance)
					}
				}
			case "withdraw":
				{
					if transaction.Amount.GreaterThan(wallet.Balance) {
						t.Error("Insufficient balance ")
					}
					err = MakeTransaction(db, c, &transaction)
					if err != nil {
						t.Errorf("Expected no error, got: %v", err)
					}
					w := &models.Wallet{}
					db.QueryRow("SELECT * FROM wallet where id_wallet=$1", transaction.WalletId).Scan(&w)
					if w.Balance.Equal(wallet.Balance) {
						t.Errorf("Expected wallet balance to be lesser than %v ", wallet.Balance)
					}
				}
			default:
				{
					err = MakeTransaction(db, c, &transaction)
					if err.Error() != "unable to treat your request" {
						t.Errorf("Expected error to be unable to treat your request, got: %v", err)
					}
				}
			}
		})
	}
}
func TestGetWallet(t *testing.T) {
	db, err := db.Connection()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	wallet := &models.Wallet{
		IdWallet: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f"),
		Balance:  decimal.NewFromFloat(19900),
		Currency: "TND",
		UserId:   uuid.MustParse("4603c129-9f40-434d-a2b2-4b0189db0da4"),
	}
	c := &gin.Context{}
	w, err := GetBalance(db, c, wallet.UserId)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if w.IdWallet != wallet.IdWallet {
		t.Errorf("Expected wallet ID to be %v, got: %v", w.IdWallet, wallet.IdWallet)
	}
}

func TestCreateWallet(t *testing.T) {
	db, err := db.Connection()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	wallets := []models.Wallet{
		{
			Balance:  decimal.NewFromFloat(19900),
			Currency: "TND",
			UserId:   uuid.MustParse("991fe02f-5f76-42a5-89b9-65d4ccf8e22e")},
		{
			Balance:  decimal.NewFromFloat(19900),
			Currency: "TND",
			UserId:   uuid.MustParse("e9e3394c-8ea6-4980-a224-054fe7a58bf6")}}
	c := &gin.Context{}
	for _, wallet := range wallets {
		t.Run(wallet.UserId.String(), func(t *testing.T) {
			if _, errs := GetUser(db, c, wallet.UserId); errs != nil {
				t.Errorf("Expected to find a user %v but got none", wallet.UserId)
			}
			err = CreateWallet(db, c, &wallet)
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}
		})
	}
}

// func TestScenario(t *testing.T) {
// 	db, err := db.Connection()
// 	if err != nil {
// 		t.Errorf("Error connecting to the database: %v", err)
// 	}
// 	defer db.Close()
// 	user := models.Users{Name: "Escanor"}

// }
