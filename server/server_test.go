package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	router := Server()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/balance/4603c129-9f40-434d-a2b2-4b0189db0da4", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestUser(t *testing.T) {
	router := Server()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/4603c129-9f40-434d-a2b2-4b0189db0da4", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	fmt.Println(w.Body)
	fmt.Println(w.Body.ReadString(','))
}

func TestUserInvalid(t *testing.T) {
	router := Server()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/89db0da4", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestDeposit(t *testing.T) {
	router := Server()
	w := httptest.NewRecorder()
	reqBody := []byte(`
	{
			"type":"deposit",
			"amount":1500
	}`)
	req, _ := http.NewRequest("POST", "/deposit/1fdec559-e095-45d0-868f-93d83b2dfb4f", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestDepositInvalid(t *testing.T) {
	router := Server()

	w := httptest.NewRecorder()
	reqBody := []byte(`
	{
			"type":"deposit",
			"amount":1500
	}`)
	req, _ := http.NewRequest("POST", "/deposit/1fdec553b2dfb4f", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestWithdraw(t *testing.T) {
	router := Server()
	w := httptest.NewRecorder()
	reqBody := []byte(`
	{
			"type":"withdraw",
			"amount":200
	}`)
	req, _ := http.NewRequest("POST", "/withdraw/1fdec559-e095-45d0-868f-93d83b2dfb4f", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestWithdrawInvalid(t *testing.T) {
	router := Server()
	w := httptest.NewRecorder()
	reqBody := []byte(`
	{
			"type":"withdraw",
			"amount":200
	}`)
	req, _ := http.NewRequest("POST", "/withdraw/1fdec559b4f", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestWallet(t *testing.T) {
	router := Server()

	recorder := httptest.NewRecorder()
	reqBody := []byte(`
	{
    "balance":2000,
    "currency":"TND"
}`)

	req, _ := http.NewRequest("POST", "/wallet/ecc62109-2815-4578-85cd-6369af0b7946", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 200, recorder.Code)

}
