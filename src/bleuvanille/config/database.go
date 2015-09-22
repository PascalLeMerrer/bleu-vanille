package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // loaded for it's side effects (defines the driver to be used)
)

var db *sql.DB

// DatabaseInit opens a connection to postgres
func DatabaseInit() {
	var err error
	connexionString := fmt.Sprintf("port=%d user=%s password=%s dbname=%s sslmode=disable", DatabasePort, DatabaseUser, DatabasePassword, DatabaseName)
	db, err = sql.Open("postgres", connexionString)
	if err != nil {
		log.Fatal(err)
	}

	createContactsTable()
	createUsersTable()
	createSessionsTable()
}

// Create Table contacts if not exists
func createContactsTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS contacts(email varchar(50) NOT NULL, created_at timestamp default NULL, constraint pk_contacts primary key(email))")

	if err != nil {
		log.Fatal(err)
	}
}

// Create Table users if not exists
func createUsersTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users(id varchar(50), email varchar(50) NOT NULL UNIQUE, firstname varchar(50), lastname varchar(50), hash varchar(100) NOT NULL, isadmin boolean, resettoken varchar(255), created_at timestamp default NULL, constraint pk_users primary key(id))")

	if err != nil {
		log.Fatal(err)
	}
}

// Create sessions table if not exists
func createSessionsTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS sessions(id varchar(50), user_id varchar(50) references users(id) ON DELETE CASCADE, is_admin boolean, expires_at timestamp default NULL, constraint pk_sessions primary key(id))")

	if err != nil {
		log.Fatal(err)
	}
}

// Db returns the database object
func Db() *sql.DB {
	if db == nil {
		DatabaseInit()
	}
	return db
}
