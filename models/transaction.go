package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	IdTransaction uuid.UUID       `JSON:"id_transaction"`
	Type          string          `JSON:"type"`
	Amount        decimal.Decimal `JSON:"amount" binding:"required"`
	Date          string          `JSON:"date"`
	WalletId      uuid.UUID       `JSON:"wallet_id"`
}
