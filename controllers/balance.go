package api

import (
	"database/sql"
	"fmt"
	models "microservice/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// deposit money
func Deposit(db *sql.DB, c *gin.Context) error {

	fail := func(err error) error {
		return fmt.Errorf("%v", err)
	}
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var transaction models.Transaction
	transaction.WalletId, _ = uuid.Parse(c.Param("id"))
	if err = c.ShouldBind(&transaction); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}
	fmt.Println(transaction)
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

// withdraw money
func Withdraw(db *sql.DB, c *gin.Context) error {
	fail := func(err error) error {
		return fmt.Errorf("%v", err)
	}
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var transaction models.Transaction
	transaction.WalletId, _ = uuid.Parse(c.Param("id"))
	if err = c.ShouldBind(&transaction); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}
	fmt.Println(transaction)

	var enough bool
	if err = tx.QueryRowContext(c, "SELECT (balance >= $1) from wallet where id_wallet = $2",
		transaction.Amount, transaction.WalletId).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("cannot find wallet"))
		}
		return fail(err)
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

// make a transaction wether it is deposit or withdraw
func Money(db *sql.DB, c *gin.Context) error {

	fail := func(err error) error {
		return fmt.Errorf("%v", err)
	}
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var transaction models.Transaction
	transaction.WalletId, _ = uuid.Parse(c.Param("id"))
	if err = c.ShouldBind(&transaction); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
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
		}
	case "withdraw":
		{
			var enough bool
			if err = tx.QueryRowContext(c, "SELECT (balance >= $1) from wallet where id_wallet = $2",
				transaction.Amount, transaction.WalletId).Scan(&enough); err != nil {
				if err == sql.ErrNoRows {
					return fail(fmt.Errorf("cannot find wallet"))
				}
				return fail(err)
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
		}
	}
	return nil
}
