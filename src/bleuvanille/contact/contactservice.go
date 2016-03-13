package contact

// Ensures persistance of contacts in Postgresql database

import (
	"bleuvanille/config"
	"encoding/json"
	"fmt"
	"github.com/solher/arangolite"
)

// CollectionName is the name of the user collection in the database
const CollectionName = "contacts"

// init creates the contacts collections if it do not already exist
func init() {
	config.DB().Run(&arangolite.CreateCollection{Name: CollectionName})
	//TODO create a hash on email and id field
}

// Save inserts a contact into the database
func Save(contact *Contact) error {
	var resultByte []byte
	var err error

	if contact.Key == "" {
		resultByte, err = config.DB().Send("INSERT DOCUMENT", "POST", "/_api/document?collection="+CollectionName, contact)
	} else {
		query := fmt.Sprintf("/_api/document/%s/%s", CollectionName, contact.Key)
		resultByte, err = config.DB().Send("UPDATE DOCUMENT", "PUT", query, contact)
	}
	if err == nil {
		err = json.Unmarshal(resultByte, contact)
	}
	return err
}

// FindByEmail returns the contact object for a given email, if any
func FindByEmail(email string) (*Contact, error) {
	var result []Contact

	query := arangolite.NewQuery(` FOR contact IN %s FILTER contact.email == @email LIMIT 1 RETURN contact `, CollectionName)
	query.Bind("email", email)

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

// FindAll returns the list of all contacts in the database
// sort defines the sorting property name
// order must be either ASC or DESC
// offset is the start index
// limit défines the max number of results to be returned
// returns an array of contacts, the total number of contacts, an error
func FindAll(sort string, order string, offset int, limit int) ([]Contact, int, error) {
	limitString := ""
	if limit > 0 {
		limitString = fmt.Sprintf("LIMIT %d, %d", offset, limit)
	}
	queryString := "FOR contact IN %s SORT contact.%s %s %s RETURN contact"
	query := arangolite.NewQuery(queryString, CollectionName, sort, order, limitString)

	async, asyncErr := config.DB().RunAsync(query)
	if asyncErr != nil {
		return nil, 0, asyncErr
	}

	contacts := []Contact{}
	decoder := json.NewDecoder(async.Buffer())

	for async.HasMore() {
		batch := []Contact{}
		decoder.Decode(&batch)
		contacts = append(contacts, batch...)
	}

	return contacts, len(contacts), nil
}

// Remove removes the entry for a given email
func Remove(email string) error {
	query := arangolite.NewQuery(`FOR contact IN %s FILTER contact.email==@email REMOVE contact IN %s`, CollectionName, CollectionName)
	query.Bind("email", email)
	_, err := config.DB().Run(query)
	return err
}
