package server

import (
	"microservice/auth"
	controller "microservice/controllers"
	db "microservice/db"

	"github.com/gin-gonic/gin"
)

func Server() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())
	db, _ := db.Connection()
	authRoutes := router.Group("/").Use(auth.AuthMiddleWare())
	//GET Method (get a user)
	//http://localhost:8080/user/52406846-1ebf-45c7-a74d-d71dcba07691
	authRoutes.GET("/user/:id", controller.GetUser(db))

	//POST Method (creating a user)
	//http://localhost:8080/user
	router.POST("/user", controller.CreateUser(db))
	//POST Method (creating a wallet)
	//http://localhost/wallet
	router.POST("/wallet/:id", controller.CreateWallet(db))

	//GET Method (view a certain person's wallet)
	//http://localhost/balance/52406846-1ebf-45c7-a74d-d71dcba07691
	authRoutes.GET("/balance/:id", controller.GetWallet(db))

	//POST Method (deposit money)
	//http://localhost/deposit/b67d93a6-465f-4e2d-ad95-8f6283397cae
	authRoutes.POST("/deposit/:id", controller.Deposit(db))

	//POST Method (withdraw money)
	//http://localhost/withdraw/b67d93a6-465f-4e2d-ad95-8f6283397cae
	authRoutes.POST("/withdraw/:id", controller.Withdraw(db))

	//POST Method (make a transaction withdraw or deposit)
	//http://localhost/money/b67d93a6-465f-4e2d-ad95-8f6283397cae
	authRoutes.POST("/money/:id", controller.MakeTransaction(db))
	return router
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Next()
	}
}
