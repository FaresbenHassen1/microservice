package models

import (
	"github.com/google/uuid"
)

type Wallet struct {
	IdWallet    uuid.UUID `JSON:"id_wallet"`
	CreatedDate string    `JSON:"created_date"`
	Balance     float32   `json:"balance"`
	Currency    string    `JSON:"currency"`
	UserId      string    `json:"users_id"`
}
