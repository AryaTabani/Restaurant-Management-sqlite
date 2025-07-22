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
	createordersTable := `
	CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	orderdate DATETIME,
	createdat DATETIME,
    updatedat DATETIME,
	orderid TEXT,
	tableid TEXT
	)
	`
	_, err = DB.Exec(createordersTable)

	if err != nil {
		panic("Could not create tables table.")
	}
	createmenusTable := `
	CREATE TABLE IF NOT EXISTS menus (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	category TEXT,
	startdate DATETIME,
	enddate DATETIME,
	createdat DATETIME,
    updatedat DATETIME,
	menuid TEXT
	)
	`
	_, err = DB.Exec(createmenusTable)

	if err != nil {
		panic("Could not create tables table.")
	}
	createinvoicesTable := `
	CREATE TABLE IF NOT EXISTS invoices (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	invoiceid TEXT,
	orderid TEXT,
	paymentmethod TEXT,
	paymentstatus TEXT,
	paymentduedate DATETIME,
	createdat DATETIME,
    updatedat DATETIME
	)
	`
	_, err = DB.Exec(createinvoicesTable)

	if err != nil {
		panic("Could not create invoices table.")
	}
	createorderitemsTable := `
	CREATE TABLE IF NOT EXISTS orderitems (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	quantity INTEGER,
	unitprice FLOAT,
	createdat DATETIME,
    updatedat DATETIME,
	foodid TEXT,
	orderitemid TEXT,
	orderid TEXT
	)
	`
	_, err = DB.Exec(createorderitemsTable)

	if err != nil {
		panic("Could not create invoices table.")
	}
}
