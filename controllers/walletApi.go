package api

import (
	"database/sql"
	models "microservice/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// uview a certain person's wallet
func GetWallet(db *sql.DB, c *gin.Context) {
	var wallet models.Wallet
	wallet.IdWallet, _ = uuid.Parse(c.Param("id"))
	err := db.QueryRow(`SELECT id_wallet, created_date, balance, currency FROM wallet WHERE users_id=$1`, wallet.IdWallet).Scan(&wallet.IdWallet, &wallet.CreatedDate, &wallet.Balance, &wallet.Currency)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(200, wallet)
}
