package config

import (
	"fmt"

	log "bleuvanille/log"

	ara "github.com/diegogub/aranGO"
)

var db *ara.Database

const (
	COLNAME_SESSIONS string = "sessions"
	COLNAME_USERS    string = "users"
	COLNAME_CONTACTS string = "contacts"
	COLNAME_ETABLES  string = "eatables"

	EDGENAME_EATABLE_PARENT string = "eatable_parent"

	GRAPHNAME_EATABLE_PARENT string = "graph_eatable_parent"
)

// DatabaseInit opens a connection to postgres
func DatabaseInit() {

	connexionString := fmt.Sprintf("http://localhost:%d", DatabasePort)

	s, err := ara.Connect(connexionString, DatabaseUser, DatabasePassword, true)
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
		log.Info(nil, "Database : Create the collection "+COLNAME_CONTACTS)

		// CollectionOptions has much more options, here we just define name , sync
		contacts := ara.NewCollectionOptions(COLNAME_CONTACTS, false)
		err := db.CreateCollection(contacts)

		if err != nil {
			log.Fatal(err)
		}

		Db().Col(COLNAME_CONTACTS).CreateHash(true, "email")
	}

	if !db.ColExist(COLNAME_USERS) {
		log.Info(nil, "Database : Create the collection "+COLNAME_USERS)

		// CollectionOptions has much more options, here we just define name , sync
		contacts := ara.NewCollectionOptions(COLNAME_USERS, false)
		err := db.CreateCollection(contacts)

		if err != nil {
			log.Fatal(err)
		}

		Db().Col(COLNAME_USERS).CreateHash(true, "email")
	}

	if !db.ColExist(COLNAME_SESSIONS) {
		log.Info(nil, "Database : Create the collection "+COLNAME_SESSIONS)

		// CollectionOptions has much more options, here we just define name , sync
		sessions := ara.NewCollectionOptions(COLNAME_SESSIONS, false)
		err := db.CreateCollection(sessions)

		if err != nil {
			log.Fatal(err)
		}
	}

	if !db.ColExist(COLNAME_ETABLES) {
		log.Info(nil, "Database : Create the collection "+COLNAME_ETABLES)

		etables := ara.NewCollectionOptions(COLNAME_ETABLES, false)
		err := db.CreateCollection(etables)

		if err != nil {
			log.Fatal(err)
		}
	}

	if db.Graph(GRAPHNAME_EATABLE_PARENT) == nil {

		edgeDefinition := ara.NewEdgeDefinition(EDGENAME_EATABLE_PARENT, []string{COLNAME_ETABLES}, []string{COLNAME_ETABLES})

		edgeDefinitionList := make([]ara.EdgeDefinition, 0)

		edgeDefinitionList = append(edgeDefinitionList, *edgeDefinition)

		_, err := db.CreateGraph(GRAPHNAME_EATABLE_PARENT, edgeDefinitionList)

		if err != nil {
			log.Fatal("Database : error when creating the graph " + GRAPHNAME_EATABLE_PARENT + " : " + err.Error())
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
