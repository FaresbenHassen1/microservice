package api

import (
	"database/sql"
	models "microservice/models"
	"microservice/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// user detail
func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.Users
		user.Id, _ = uuid.Parse(c.Param("id"))
		u, err := service.GetUser(db, c, user.Id)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		c.JSON(200, u)
	}
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.Users
		if err := ctx.ShouldBind(&user); err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		err := service.CreateUser(db, ctx, user.Name)
		if err != nil {
			ctx.AbortWithStatusJSON(400, err.Error())
		}
		ctx.JSON(200, user)
	}

}
