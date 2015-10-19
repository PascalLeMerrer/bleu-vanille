package config

import (
	"fmt"
	"log"

	ara "github.com/diegogub/aranGO"
)

var db *ara.Database

const COLNAME_SESSIONS string = "sessions"
const COLNAME_USERS string = "users"
const COLNAME_CONTACTS string = "contacts"

// DatabaseInit opens a connection to postgres
func DatabaseInit() {

	connexionString := fmt.Sprintf("http://localhost:%d", DatabasePort)

	s, err := ara.Connect(connexionString, DatabaseUser, DatabasePassword, false)
	if err != nil {
		log.Fatal(err)
	}

	db = s.DB(DatabaseName)

	//Create the database if necessary
	if db == nil {
		user := ara.User{Username: DatabaseUser, Password: DatabasePassword}
		users := []ara.User{user}
		err = s.CreateDB(DatabaseName, users)

		if err != nil {
			log.Fatal(err)
		}

		db = s.DB(DatabaseName)
	}

	//Create the tables if necessary
	createTables()
}

// Create Table contacts if not exists
func createTables() {
	if !db.ColExist(COLNAME_CONTACTS) {
		// CollectionOptions has much more options, here we just define name , sync
		contacts := ara.NewCollectionOptions(COLNAME_CONTACTS, true)
		err := db.CreateCollection(contacts)

		if err != nil {
			log.Fatal(err)
		}
	}

	if !db.ColExist(COLNAME_USERS) {
		// CollectionOptions has much more options, here we just define name , sync
		contacts := ara.NewCollectionOptions(COLNAME_USERS, true)
		err := db.CreateCollection(contacts)

		if err != nil {
			log.Fatal(err)
		}
	}

	if !db.ColExist(COLNAME_SESSIONS) {
		// CollectionOptions has much more options, here we just define name , sync
		contacts := ara.NewCollectionOptions(COLNAME_SESSIONS, true)
		err := db.CreateCollection(contacts)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetCollection(m ara.Modeler) *ara.Collection {
	return Db().Col(m.GetCollection())
}

// Db returns the database object
func Db() *ara.Database {
	if db == nil {
		DatabaseInit()
	}
	return db
}

func Context() *ara.Context {
	ctx, err := ara.NewContext(Db())

	if err != nil {
		log.Printf(err.Error())
	}

	return ctx
}