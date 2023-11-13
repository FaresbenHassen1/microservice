package api

import (
	"database/sql"
	models "microservice/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// user detail
func GetUser(db *sql.DB, c *gin.Context) {
	var user models.Users
	user.Id, _ = uuid.Parse(c.Param("id"))

	err := db.QueryRow(`SELECT id_user, name FROM users WHERE id_user=$1`, user.Id).Scan(&user.Id, &user.Name)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(200, user)

}
