package session

import (
	"bleuvanille/config"
	"encoding/json"
	"fmt"
	"github.com/solher/arangolite"
)

// CollectionName is the name of the user collection in the database
const CollectionName = "sessions"

// init creates the contacts collections if it do not already exist
func init() {
	config.DB().Run(&arangolite.CreateCollection{Name: CollectionName})
	//TODO create a hash on id field
}

// Save inserts a session into the database
// TODO we should renew the id is the admin privilege is given
// See https://www.owasp.org/index.php/Session_Management_Cheat_Sheet
func Save(sess *Session) error {
	var resultByte []byte
	var err error

	if sess.Key == "" {
		resultByte, err = config.DB().Send("INSERT DOCUMENT", "POST", "/_api/document?collection="+CollectionName, sess)
	} else {
		query := fmt.Sprintf("/_api/document/%s/%s", CollectionName, sess.Key)
		resultByte, err = config.DB().Send("UPDATE DOCUMENT", "PUT", query, sess)
	}
	if err == nil {
		err = json.Unmarshal(resultByte, sess)
	}
	return err
}

// FindByID returns the sess object for a given ID, if any,
// or nil if the session does not exist in the database
func FindByID(ID string) (*Session, error) {
	var result []Session

	query := arangolite.NewQuery(` FOR session IN %s FILTER session.SessionID == @id LIMIT 1 RETURN session `, CollectionName)
	query.Bind("id", ID)

	rawResult, err := config.DB().Run(query)
	if err != nil {
		return nil, err
	}

	marshallErr := json.Unmarshal(rawResult, &result)
	if marshallErr != nil {
		return nil, marshallErr
	}
	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, nil
}

// Remove removes from the database the session for the given ID
func Remove(ID string) error {
	query := arangolite.NewQuery(`FOR session IN %s FILTER session.SessionID==@id REMOVE session IN %s`, CollectionName, CollectionName)
	query.Bind("id", ID)
	_, err := config.DB().Run(query)
	return err
}
