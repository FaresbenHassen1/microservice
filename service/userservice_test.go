package service

import (
	"microservice/db"
	"microservice/models"
	"microservice/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) models.Users {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	user1 := models.Users{
		Name: util.RandomName(),
	}
	user2, err := CreateUser(db, &gin.Context{}, user1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Name, user2.Name)
	require.NotZero(t, user2.Id)
	return user2
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	db, err := db.ConnectTestDB()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	user1 := CreateRandomUser(t)
	user2, err := GetUser(db, &gin.Context{}, user1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Id, user2.Id)
	require.Equal(t, user1.Name, user2.Name)
}
