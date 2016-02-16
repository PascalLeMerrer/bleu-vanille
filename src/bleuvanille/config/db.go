package config

import (
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
		LoggerOptions(Debug, Debug, Debug).
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
