package server

import (
	controller "microservice/controllers"
	db "microservice/db"
	models "microservice/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Server() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())
	db, err := db.Connection()

	//GET Method (get a user)
	//http://localhost:8080/user/52406846-1ebf-45c7-a74d-d71dcba07691
	router.GET("/user/:id", func(c *gin.Context) {
		controller.GetUser(db, c)
	})

	//POST Method (creating a wallet)
	//http://localhost/wallet
	router.POST("/wallet", func(c *gin.Context) {
		var w models.Wallet
		if err = c.ShouldBind(&w); err != nil {
			c.AbortWithStatusJSON(400, err.Error())
		}
		id := ""
		sqlPostStatement := `INSERT INTO wallet (balance,currency) VALUES ($1,$2) RETURNING id_wallet`
		err = db.QueryRow(sqlPostStatement, w.Balance, w.Currency).Scan(&id)
		if err != nil {
			panic(err)
		}
		c.JSON(200, w)
	})

	//GET Method (view a certain person's wallet)
	//http://localhost/balance/52406846-1ebf-45c7-a74d-d71dcba07691
	router.GET("/balance/:id", func(c *gin.Context) {
		controller.GetWallet(db, c)
	})

	//POST Method (deposit money)
	//http://localhost/deposit/b67d93a6-465f-4e2d-ad95-8f6283397cae
	router.POST("/deposit/:id", func(c *gin.Context) {
		err := controller.Deposit(db, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deposit successfully"})
	})

	//POST Method (withdraw money)
	//http://localhost/withdraw/b67d93a6-465f-4e2d-ad95-8f6283397cae
	router.POST("/withdraw/:id", func(c *gin.Context) {
		err := controller.Withdraw(db, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "withdraw successfully"})
	})

	//POST Method (make a transaction withdraw or deposit)
	//http://localhost/money/b67d93a6-465f-4e2d-ad95-8f6283397cae
	router.POST("/money/:id", func(c *gin.Context) {
		err := controller.Money(db, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Transaction done successfully"})
	})
	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Next()
	}
}
