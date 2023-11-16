package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	IdWallet    uuid.UUID       `JSON:"id_wallet"`
	CreatedDate string          `JSON:"created_date"`
	Balance     decimal.Decimal `json:"balance"`
	Currency    string          `JSON:"currency"`
	UserId      uuid.UUID       `json:"users_id"`
}
