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

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`
	_, err = DB.Exec(createEventsTable)
	if err != nil {
		panic("could not create events table: " + err.Error())
	}

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registration (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id INTEGER,
	user_id INTEGER,
	FOREIGN KEY(event_id) REFERENCES events(id)
	FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createRegistrationsTable)

	if err != nil {
		panic("Could not create registration table.")
	}
}
