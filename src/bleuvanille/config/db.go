package config

import (
	"bleuvanille/log"
	"fmt"
	"github.com/solher/arangolite"
)

var dbSession *arangolite.DB

// DB returns a connection to ArangoDB
func DB() *arangolite.DB {
	if dbSession != nil {
		return dbSession
	}

	dbSession = arangolite.New().
		LoggerOptions(DbDebug(), DbDebug(), DbDebug()).
		Connect(DatabaseProtocol+"://"+DatabaseHost+":"+fmt.Sprint(DatabasePort),
		"_system",
		"root",
		DatabaseRootPassword)

	// we try to create the database; it will fail if it already exists
	dbSession.Run(&arangolite.CreateDatabase{
		Name: DatabaseName,
		Users: []map[string]interface{}{
			{"username": DatabaseUser, "passwd": DatabasePassword},
		},
	})

	dbSession.SwitchDatabase(DatabaseName).SwitchUser(DatabaseUser, DatabasePassword)
	return dbSession
}

// CreateHashIndexedCollection creates a collection
// name is the collection name
// indexes is array of hash indexes
func CreateHashIndexedCollection(name string, indexes []string) {
	rawResult, err := DB().Run(&arangolite.CreateCollection{Name: name})
	if err != nil {
		fmt.Printf("ERROR: cannot create %s collection: %s \n", name, err)
		return
	}
	log.Printf("%s collection created: %s \n", name, rawResult)
	if len(indexes) == 0 {
		return
	}
	unique := true
	sparse := false
	hashIndex := arangolite.CreateHashIndex{
		CollectionName: name,
		Unique:         &unique,
		Sparse:         &sparse,
		Fields:         indexes,
	}
	rawResult, err = DB().Run(&hashIndex)
	if err != nil {
		fmt.Printf("ERROR: cannot add index to %s collection %s \n", name, err)
		return
	}
	log.Printf("Index added to %s collection: %s \n", name, rawResult)
}
