package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	env "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "goproject"
// )

func Connection() (*sql.DB, error) {

	err := env.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	fmt.Println(os.Getenv("HOST"))
	fmt.Println(os.Getenv("DATABASE_PORT"))
	fmt.Println(os.Getenv("DATABASE_USER"))
	fmt.Println(os.Getenv("DATABASE_PASSWORD"))
	fmt.Println(os.Getenv("DATABASE_NAME"))

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Simply connecting to DataBase")
	return db, err
}
