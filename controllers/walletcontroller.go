package api

import (
	"database/sql"
	models "microservice/models"
	service "microservice/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// view a certain person's wallet
func GetWallet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var wallet models.Wallet
		wallet.UserId, _ = uuid.Parse(c.Param("id"))
		w, err := service.GetBalance(db, c, wallet.UserId)
		if err != nil {
			c.AbortWithStatusJSON(400, err.Error())
		}
		c.JSON(http.StatusOK, w)
	}
}

// create a wallet
func CreateWallet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var wallet models.Wallet
		wallet.UserId, _ = uuid.Parse(c.Param("id"))
		if err := c.ShouldBind(&wallet); err != nil {
			c.AbortWithStatusJSON(400, err.Error())
		}
		_, err := service.CreateWallet(db, c, &wallet)
		if err != nil {
			c.AbortWithStatusJSON(400, err.Error())
		}
		c.JSON(200, wallet)
	}
}
