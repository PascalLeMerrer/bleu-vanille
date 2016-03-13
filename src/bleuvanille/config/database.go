package config

import (
	"fmt"

	log "bleuvanille/log"

	ara "github.com/diegogub/aranGO"
)

var db *ara.Database

const (
	COLNAME_SESSIONS string = "sessions"
	COLNAME_ETABLES  string = "eatables"

	EDGENAME_EATABLE_PARENT string = "eatable_parent"

	GRAPHNAME_EATABLE_PARENT string = "graph_eatable_parent"
)

// DatabaseInit opens a connection to ArangoDB and creates the database if it doesn't exist
func DatabaseInit() {

	connexionString := fmt.Sprintf("http://localhost:%d", DatabasePort)

	s, err := ara.Connect(connexionString, DatabaseUser, DatabasePassword, false)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot connect to database %v", err))
	}

	availableDBs, dbListError := s.AvailableDBs()

	if dbListError != nil {
		log.Fatal(fmt.Sprintf("Cannot list available databases %v", err))
	}

	for _, dbName := range availableDBs {
		if dbName == DatabaseName {
			db = s.DB(DatabaseName)
		}
	}

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

	if !db.ColExist(COLNAME_SESSIONS) {
		log.Info(nil, "Database: Creating the collection "+COLNAME_SESSIONS)

		// CollectionOptions has much more options, here we just define name , sync
		sessions := ara.NewCollectionOptions(COLNAME_SESSIONS, false)
		err := db.CreateCollection(sessions)

		if err != nil {
			log.Fatal(err)
		}
	}
}

//GetCollection returns the collection object related to the given object modeler
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

//Context returns the context related to the DB
func Context() *ara.Context {
	ctx, err := ara.NewContext(Db())

	if err != nil {
		log.Printf(err.Error())
	}

	return ctx
}
