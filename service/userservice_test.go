package service

import (
	"microservice/db"
	"microservice/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userData struct {
	id   uuid.UUID
	name string
}

func TestGetUser(t *testing.T) {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	usersData := []userData{
		{id: uuid.MustParse("4603c129-9f40-434d-a2b2-4b0189db0da4"), name: "fares"},
		{id: uuid.MustParse("ecc62109-2815-4578-85cd-6369af0b7946"), name: "Tirion"},
		{id: uuid.MustParse("abcdef12-3456-7890-0123-456789abcdef"), name: "fares"},
	}
	// should pass 1st case no issue
	// should fails 2nd case because expected name
	// should fails 3rd case because user not found
	for _, userData := range usersData {
		t.Run(userData.name, func(t *testing.T) {
			c := &gin.Context{}
			user, err := GetUser(db, c, userData.id)
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if user.Name != userData.name {
				t.Errorf("Expected user name to be %v, got: %v", userData.name, user.Name)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	users := []models.Users{{Name: "Garrosh"}, {}}
	c := &gin.Context{}
	for _, user := range users {
		t.Run(user.Name, func(t *testing.T) {
			if len(user.Name) < 1 {
				t.Errorf("Expected user to have a name")
			}
			err = CreateUser(db, c, user.Name)
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
