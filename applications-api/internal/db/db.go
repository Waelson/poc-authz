package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	ID     string `json:"application_id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "applications.db")
	if err != nil {
		panic(err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS tb_application (
        application_id VARCHAR(150) NOT NULL,
		user_id INT NOT NULL,
        name VARCHAR(50) NOT NULL
    );`
	_, err = DB.Exec(createTable)
	if err != nil {
		panic(err)
	}
}
