package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Relationship struct {
	RelationType string `json:"relation_type"`
	Resource     string `json:"resource"`
}

type Operator struct {
	RootID       int            `json:"root_id"`
	OperatorID   int            `json:"operator_id"`
	RoleID       string         `json:"role_id"`
	Relationship []Relationship `json:"relationship"`
}

type OperatorResponse struct {
	RootID     int    `json:"root_id"`
	OperatorID int    `json:"operator_id"`
	RoleID     string `json:"role_id"`
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "operators.db")
	if err != nil {
		panic(err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS tb_operator (
        root_id INTEGER,
        operator_id INTEGER PRIMARY KEY,
        role_id VARCHAR(200),
        relationship TEXT
    );`
	_, err = DB.Exec(createTable)
	if err != nil {
		panic(err)
	}
}
