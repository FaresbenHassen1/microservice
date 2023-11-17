package service

import (
	"database/sql"
	"fmt"
	"microservice/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUser(db *sql.DB, c *gin.Context, id uuid.UUID) (*models.Users, error) {
	user := &models.Users{}
	user.Id = id
	err := db.QueryRow(`SELECT id_user, name FROM users WHERE id_user=$1`, user.Id).Scan(&user.Id, &user.Name)
	if err != nil {
		return user, fmt.Errorf("this person does not exist")
	}
	return user, err
}

func CreateUser(db *sql.DB, c *gin.Context, name string) error {
	var user models.Users
	if 1 > len(name) {
		return fmt.Errorf("unable to create a user with an empty name")
	}
	err := db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id_user", name).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
