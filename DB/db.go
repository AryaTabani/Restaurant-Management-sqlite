package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to database")
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createtables()
}

func createtables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    firstname TEXT,
    lastname TEXT,
    password TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    avatar TEXT,
    phone TEXT,
    createdat DATETIME,
    updatedat DATETIME,
    userid TEXT
)
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("could not create users table: " + err.Error())
	}



	createFoodsTable := `
	CREATE TABLE IF NOT EXISTS foods (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	price TEXT,
	foodimage TEXT,
	createdat DATETIME,
    updatedat DATETIME,
	foodid TEXT,
	menuid TEXT
	)
	`
	_, err = DB.Exec(createFoodsTable)

	

	if err != nil {
		panic("Could not create foods table.")
	}

	createtablesTable := `
	CREATE TABLE IF NOT EXISTS tables (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	numberofguests INTEGER,
	tablenumber INTEGER,
	createdat DATETIME,
    updatedat DATETIME,
	tableid TEXT
	)
	`
	_, err = DB.Exec(createtablesTable)

	

	if err != nil {
		panic("Could not create tables table.")
	}
}
