package models

import (
	"github.com/google/uuid"
)

type Transaction struct {
	IdTransaction uuid.UUID `JSON:"id_transaction"`
	Type          string    `JSON:"type"`
	Amount        float32   `JSON:"amount"`
	Date          string    `JSON:"date"`
	WalletId      uuid.UUID `JSON:"wallet_id"`
}
