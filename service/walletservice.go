package service

import (
	"database/sql"
	"fmt"
	"microservice/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Withdraw(db *sql.DB, c *gin.Context, t *models.Transaction) error {
	fail := func(err error) error {
		return fmt.Errorf("%v", err)
	}
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var transaction = *t
	if transaction.Amount.IsNegative() {
		return fmt.Errorf("the amount of money to withdraw is negative, cannot proceed with transaction")
	}
	if transaction.Type != "withdraw" {
		return fmt.Errorf("invalid transaction type")
	}
	var enough bool
	if err = tx.QueryRowContext(c, "SELECT (balance >= $1) from wallet where id_wallet = $2",
		transaction.Amount, transaction.WalletId).Scan(&enough); err != nil {
		return fail(fmt.Errorf("cannot find wallet"))
	}
	if !enough {
		return fail(fmt.Errorf("not enough money, your broke as hell"))
	}
	_, err = tx.ExecContext(c, "INSERT INTO transaction (type,amount,wallet_id) VALUES ($1,$2,$3)", transaction.Type, transaction.Amount, transaction.WalletId)
	if err != nil {
		return fail(err)
	}
	_, err = tx.ExecContext(c, "UPDATE wallet SET balance = balance - $1 WHERE id_wallet=$2", transaction.Amount, transaction.WalletId)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return nil
}

func Deposit(db *sql.DB, c *gin.Context, t *models.Transaction) error {
	fail := func(err error) error {
		return fmt.Errorf("%v", err)
	}
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var transaction = *t
	if transaction.Type == "deposit" {
		if transaction.Amount.IsNegative() {
			return fmt.Errorf("the amount of money to deposit is negative, cannot proceed with transaction")
		}
		_, err = tx.ExecContext(c, "INSERT INTO transaction (type,amount,wallet_id) VALUES ($1,$2,$3)", transaction.Type, transaction.Amount, transaction.WalletId)
		if err != nil {
			return fail(err)
		}
		_, err = tx.ExecContext(c, "UPDATE wallet SET balance = balance + $1 WHERE id_wallet=$2", transaction.Amount, transaction.WalletId)
		if err != nil {
			return fail(err)
		}
		if err = tx.Commit(); err != nil {
			return fail(err)
		}
		return nil
	} else {
		return fmt.Errorf("invalid transaction type")
	}

}

func MakeTransaction(db *sql.DB, c *gin.Context, t *models.Transaction) error {
	fail := func(err error) error {
		return fmt.Errorf("%v", err)
	}
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var transaction = *t
	if transaction.Amount.IsNegative() {
		return fmt.Errorf("the amount of money to deposit is negative, cannot proceed with transaction")
	}
	switch transaction.Type {
	case "deposit":
		{
			_, err = tx.ExecContext(c, "INSERT INTO transaction (type,amount,wallet_id) VALUES ($1,$2,$3)", transaction.Type, transaction.Amount, transaction.WalletId)
			if err != nil {
				return fail(err)
			}
			_, err = tx.ExecContext(c, "UPDATE wallet SET balance = balance + $1 WHERE id_wallet=$2", transaction.Amount, transaction.WalletId)
			if err != nil {
				return fail(err)
			}
			if err = tx.Commit(); err != nil {
				return fail(err)
			}
			return nil
		}
	case "withdraw":
		{
			var enough bool
			if err = tx.QueryRowContext(c, "SELECT (balance >= $1) from wallet where id_wallet = $2",
				transaction.Amount, transaction.WalletId).Scan(&enough); err != nil {
				return fail(fmt.Errorf("cannot find wallet"))
			}
			if !enough {
				return fail(fmt.Errorf("not enough money, your broke as hell"))
			}
			_, err = tx.ExecContext(c, "INSERT INTO transaction (type,amount,wallet_id) VALUES ($1,$2,$3)", transaction.Type, transaction.Amount, transaction.WalletId)
			if err != nil {
				return fail(err)
			}
			_, err = tx.ExecContext(c, "UPDATE wallet SET balance = balance - $1 WHERE id_wallet=$2", transaction.Amount, transaction.WalletId)
			if err != nil {
				return fail(err)
			}
			if err = tx.Commit(); err != nil {
				return fail(err)
			}
			return nil
		}
	default:
		if transaction.Type != "deposit" && transaction.Type != "withdraw" {
			return fail(fmt.Errorf("unable to treat your request"))
		}
	}
	return nil

}

func GetBalance(db *sql.DB, c *gin.Context, id uuid.UUID) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	wallet.UserId = id
	err := db.QueryRow(`SELECT id_wallet, created_date, balance, currency FROM wallet WHERE users_id=$1`, wallet.UserId).Scan(&wallet.IdWallet, &wallet.CreatedDate, &wallet.Balance, &wallet.Currency)
	if err != nil {
		return wallet, fmt.Errorf("cannot find wallet")
	}
	return wallet, err
}

func CreateWallet(db *sql.DB, c *gin.Context, w *models.Wallet) error {
	var wallet = *w
	if _, err := GetUser(db, c, wallet.UserId); err != nil {
		return fmt.Errorf("this person does not exist")
	} else if _, err = GetBalance(db, c, w.UserId); err == nil {
		return fmt.Errorf("this person already posses a wallet, you cannot add more wallet")
	} else {
		err := db.QueryRow("INSERT INTO wallet (balance,currency,users_id) VALUES ($1,$2,$3) RETURNING id_wallet", wallet.Balance, wallet.Currency, wallet.UserId).Scan(&wallet.IdWallet)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		return nil
	}

}
