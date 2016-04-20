package user

// Ensures persistance of user accounts in Postgresql database

import (
	"bleuvanille/config"
	"encoding/json"
	"fmt"

	"github.com/solher/arangolite"
)

// CollectionName is the name of the user collection in the database
const CollectionName = "users"

// init creates the users collections if it do not already exist
func init() {
	indexFields := []string{"email"}
	config.CreateHashIndexedCollection(CollectionName, indexFields)
}

// Save inserts a user into the database
func Save(user *User) error {
	var resultByte []byte
	var err error

	if user.Key == "" {
		resultByte, err = config.DB().Send("INSERT DOCUMENT", "POST", "/_api/document?collection="+CollectionName, user)
	} else {
		query := fmt.Sprintf("/_api/document/%s/%s", CollectionName, user.Key)
		resultByte, err = config.DB().Send("UPDATE DOCUMENT", "PUT", query, user)
	}
	if err == nil {
		err = json.Unmarshal(resultByte, user)
	}
	return err
}

// SavePassword updates the password of a given user into the database
func SavePassword(user *User, newPassword string) error {
	user.Hash = newPassword
	return Save(user)
}

// FindByEmail returns the user object for a given email, if any
func FindByEmail(email string) (*User, error) {
	return FindBy("email", email)
}

// FindByID returns the user object for a given user ID, if any
func FindByID(ID string) (*User, error) {
	return FindBy("id", ID)
}

// FindByName returns the user object for a given name, if any
// TODO: what if two users have the same name?
func FindByName(name string) (*User, error) {
	return FindBy("name", name)
}

// FindBy returns an user matching the given property value, or nil if not user matches
func FindBy(name, value string) (*User, error) {

	query := arangolite.NewQuery(` FOR u IN %s FILTER u.@name == @value LIMIT 1 RETURN u `, CollectionName)
	query.Bind("name", name)
	query.Bind("value", value)

	return executeReadingQuery(query)
}

// Executes a given query that is expected to return a single user
func executeReadingQuery(query *arangolite.Query) (*User, error) {
	var result []User

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

// FindAll returns the list of all users in the database
// sort defines the sorting property name
// order must be either ASC or DESC
// offset is the start index
// limit defines the max number of results to be returned
// returns an array of user, the total number of users, an error
func FindAll(sort string, order string, offset int, limit int) ([]User, int, error) {

	limitString := ""
	if limit > 0 {
		limitString = fmt.Sprintf("LIMIT %d, %d", offset, limit)
	}
	queryString := "FOR u IN %s SORT u.%s %s %s RETURN u"
	query := arangolite.NewQuery(queryString, CollectionName, sort, order, limitString)

	async, asyncErr := config.DB().RunAsync(query)
	if asyncErr != nil {
		return nil, 0, asyncErr
	}

	users := []User{}
	decoder := json.NewDecoder(async.Buffer())

	for async.HasMore() {
		batch := []User{}
		decoder.Decode(&batch)
		users = append(users, batch...)
	}

	return users, len(users), nil

}

// Remove deletes the user for the given key from the database
func Remove(key string) error {
	query := arangolite.NewQuery(`REMOVE @key IN %s`, CollectionName)
	query.Bind("key", key)
	_, err := config.DB().Run(query)
	return err
}
