package api

import (
	"database/sql"
	models "microservice/models"
	service "microservice/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// deposit money
func Deposit(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var transaction models.Transaction
		transaction.WalletId, _ = uuid.Parse(ctx.Param("id"))
		if err := ctx.BindJSON(&transaction); err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		if _, _, err := service.Deposit(db, ctx, &transaction); err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		ctx.String(http.StatusOK, "", transaction)
	}
}

// withdraw money
func Withdraw(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var transaction models.Transaction
		transaction.WalletId, _ = uuid.Parse(ctx.Param("id"))
		if err := ctx.BindJSON(&transaction); err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		if _, _, err := service.Withdraw(db, ctx, &transaction); err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		ctx.String(http.StatusOK, "", transaction)
	}
}

// make a transaction wether it is deposit or withdraw
func MakeTransaction(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var transaction models.Transaction
		transaction.WalletId, _ = uuid.Parse(ctx.Param("id"))
		if err := ctx.ShouldBind(&transaction); err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		if err := service.MakeTransaction(db, ctx, &transaction); err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		ctx.String(http.StatusOK, "", transaction)
	}
}
