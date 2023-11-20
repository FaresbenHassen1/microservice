package service

import (
	"microservice/db"
	"microservice/models"
	"microservice/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// func TestMakeTransaction(t *testing.T) {
// 	db, err := db.ConnectTestDB()
// 	if err != nil {
// 		t.Errorf("Error connecting to the database: %v", err)
// 	}
// 	defer db.Close()
// 	c := &gin.Context{}
// 	transactions := []models.Transaction{
// 		{Type: "withdraw",
// 			Amount:   decimal.NewFromFloat(1000),
// 			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
// 		{Type: "deposit",
// 			Amount:   decimal.NewFromFloat(1000),
// 			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
// 		{Type: "z",
// 			Amount:   decimal.NewFromFloat(1000),
// 			WalletId: uuid.MustParse("1fdec559-e095-45d0-868f-93d83b2dfb4f")},
// 	}
// 	for _, transaction := range transactions {
// 		t.Run(transaction.Type, func(t *testing.T) {
// 			if transaction.Amount.IsNegative() {
// 				t.Errorf("the amount of money to withdraw is negative, cannot proceed with transaction")
// 			}
// 			wallet := &models.Wallet{}
// 			err = db.QueryRow("SELECT balance FROM wallet WHERE id_wallet=$1", transaction.WalletId).Scan(&wallet.Balance)
// 			if err != nil {
// 				t.Errorf("Error finding wallet: %v", err)
// 			}
// 			switch transaction.Type {
// 			case "deposit":
// 				{
// 					err = MakeTransaction(db, c, &transaction)
// 					if err != nil {
// 						t.Errorf("Expected no error, got: %v", err)
// 					}
// 					w := &models.Wallet{}
// 					db.QueryRow("SELECT * FROM wallet where id_wallet=$1", transaction.WalletId).Scan(&w)
// 					if w.Balance.Equal(wallet.Balance) {
// 						t.Errorf("Expected wallet balance to be higher than %v ", wallet.Balance)
// 					}
// 				}
// 			case "withdraw":
// 				{
// 					if transaction.Amount.GreaterThan(wallet.Balance) {
// 						t.Error("Insufficient balance ")
// 					}
// 					err = MakeTransaction(db, c, &transaction)
// 					if err != nil {
// 						t.Errorf("Expected no error, got: %v", err)
// 					}
// 					w := &models.Wallet{}
// 					db.QueryRow("SELECT * FROM wallet where id_wallet=$1", transaction.WalletId).Scan(&w)
// 					if w.Balance.Equal(wallet.Balance) {
// 						t.Errorf("Expected wallet balance to be lesser than %v ", wallet.Balance)
// 					}
// 				}
// 			default:
// 				{
// 					err = MakeTransaction(db, c, &transaction)
// 					if err.Error() != "unable to treat your request" {
// 						t.Errorf("Expected error to be unable to treat your request, got: %v", err)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

func CreateRandomWallet(t *testing.T) models.Wallet {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	user := CreateRandomUser(t)
	wallet := models.Wallet{
		Balance:  decimal.NewFromInt(util.RandomMoney()),
		Currency: util.RandomCurrency(),
		UserId:   user.Id,
	}
	wallet2, err := CreateWallet(db, &gin.Context{}, &wallet)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)
	require.NotZero(t, wallet2.IdWallet)
	require.Equal(t, wallet.Balance, wallet2.Balance)
	require.Equal(t, wallet.Currency, wallet2.Currency)
	require.Equal(t, wallet.UserId, wallet2.UserId)
	require.NotNil(t, wallet2.CreatedDate)
	return wallet2
}

func TestCreateWallet(t *testing.T) {
	CreateRandomWallet(t)
}

func TestGetWallet(t *testing.T) {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	wallet := CreateRandomWallet(t)
	wallet2, err := GetBalance(db, &gin.Context{}, wallet.UserId)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)
	require.Equal(t, wallet.IdWallet, wallet2.IdWallet)
	require.Equal(t, wallet.Balance, wallet2.Balance)
	require.Equal(t, wallet.Currency, wallet2.Currency)
	require.Equal(t, wallet.UserId, wallet2.UserId)
}

func TestDeposit(t *testing.T) {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	wallet := CreateRandomWallet(t)
	transaction := models.Transaction{
		Amount:   decimal.NewFromInt(util.RandomAmount()),
		Type:     "deposit",
		WalletId: wallet.IdWallet,
	}
	transaction2, wallet2, err := Deposit(db, &gin.Context{}, &transaction)
	require.NoError(t, err)
	require.NotEmpty(t, transaction2)
	require.Equal(t, transaction.Amount, transaction2.Amount)
	require.Equal(t, transaction.Type, transaction2.Type)
	require.Equal(t, transaction.WalletId, transaction2.WalletId)
	require.NotNil(t, transaction2.IdTransaction)
	require.Equal(t, wallet2.Balance, decimal.Sum(wallet.Balance, transaction.Amount))
}

func TestWithdraw(t *testing.T) {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	wallet := CreateRandomWallet(t)
	transaction := models.Transaction{
		Amount:   decimal.NewFromInt(util.RandomAmount()),
		Type:     "withdraw",
		WalletId: wallet.IdWallet,
	}

	transaction2, wallet2, err := Withdraw(db, &gin.Context{}, &transaction)
	require.NoError(t, err)
	require.NotEmpty(t, transaction2)
	require.Equal(t, transaction.Amount, transaction2.Amount)
	require.Equal(t, transaction.Type, transaction2.Type)
	require.Equal(t, transaction.WalletId, transaction2.WalletId)
	require.NotNil(t, transaction2.IdTransaction)
	require.Equal(t, wallet2.Balance, wallet.Balance.Sub(transaction2.Amount))
}
