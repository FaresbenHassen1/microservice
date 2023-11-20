package server

import (
	"bytes"
	"encoding/json"
	controller "microservice/controllers"
	db "microservice/db"
	"microservice/models"
	"microservice/service"
	"microservice/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func CreateRandomUser(t *testing.T) models.Users {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	user1 := models.Users{
		Name: util.RandomName(),
	}
	user2, _ := service.CreateUser(db, &gin.Context{}, user1.Name)
	return user2
}

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
	wallet2, _ := service.CreateWallet(db, &gin.Context{}, &wallet)
	return wallet2
}
func TestBalance(t *testing.T) {
	router := gin.Default()
	db, _ := db.ConnectTestDB()
	router.GET("/balance/:id", controller.GetWallet(db))
	wallet := CreateRandomWallet(t)
	jsonValue, _ := json.Marshal(wallet)
	req, _ := http.NewRequest("GET", "/balance/"+wallet.UserId.String(), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	reqNotFound, _ := http.NewRequest("GET", "/balance/10", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, reqNotFound)
	assert.Equal(t, 400, w.Code)
}

func TestUser(t *testing.T) {
	router := gin.Default()
	db, _ := db.ConnectTestDB()
	router.GET("/user/:id", controller.GetUser(db))
	user := CreateRandomUser(t)
	jsonValue, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/"+user.Id.String(), bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	reqNotFound, _ := http.NewRequest("GET", "/user/10", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, reqNotFound)
	assert.Equal(t, 400, w.Code)
}

func TestDeposit(t *testing.T) {
	router := gin.Default()
	db, _ := db.ConnectTestDB()
	router.POST("/deposit/:id", controller.Deposit(db))
	wallet := CreateRandomWallet(t)
	transaction := models.Transaction{
		Amount:   decimal.NewFromInt(util.RandomAmount()),
		Type:     "deposit",
		WalletId: wallet.IdWallet,
	}
	jsonValue, _ := json.Marshal(transaction)
	req, _ := http.NewRequest("POST", "/deposit/"+transaction.WalletId.String(), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestWithdraw(t *testing.T) {
	router := gin.Default()
	db, _ := db.ConnectTestDB()
	router.POST("/withdraw/:id", controller.Withdraw(db))
	wallet := CreateRandomWallet(t)
	transaction := models.Transaction{
		Amount:   decimal.NewFromInt(util.RandomAmount()),
		Type:     "withdraw",
		WalletId: wallet.IdWallet,
	}
	jsonValue, _ := json.Marshal(transaction)
	req, _ := http.NewRequest("POST", "/withdraw/"+transaction.WalletId.String(), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
