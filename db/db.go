package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/f-chilmi/just-text-go/helpers"
	"github.com/joho/godotenv"
)

func CreateConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	helpers.CheckError("Error loading env files (create connection)", err)

	// initialize db credential
	DbHost := os.Getenv("HOST")
	DbPort := os.Getenv("PORT")
	DbUser := os.Getenv("USER")
	DbPassword := os.Getenv("PASSWORD")
	DbName := os.Getenv("DBNAME")

	DbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	// open the connection
	db, err := sql.Open("postgres", DbUrl)
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	// return the connection
	return db
}
