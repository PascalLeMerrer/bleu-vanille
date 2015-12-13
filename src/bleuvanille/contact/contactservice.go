package contact

// Ensures persistance of contacts in Postgresql database

import (
	"bleuvanille/config"
	"errors"
	"fmt"
	ara "github.com/diegogub/aranGO"
	"math"
	"strconv"
)

// Save inserts a contact into the database
func Save(contact *Contact) error {
	errorMap := config.Context().Save(contact)
	if value, ok := errorMap["error"]; ok {
		return errors.New(value)
	}
	return nil
}

// LoadByEmail returns the contact object for a given email, if any
func LoadByEmail(email string) (*Contact, error) {
	var result Contact

	col := config.GetCollection(&result)
	cursor, err := col.Example(map[string]interface{}{"email": email}, 0, 1)
	if err != nil {
		return nil, err
	}
	if cursor.Result != nil && len(cursor.Result) > 0 {
		cursor.FetchOne(&result)
		return &result, nil
	}
	return nil, nil
}

// LoadAll returns the list of all contacts in the database
// sort defines the sorting property name
// order must be either ASC or DESC
// offset is the start index
// limit défines the max number of results to be returned
// returns an array of contacts, the total number of contacts, an error
func LoadAll(sort string, order string, offset int, limit int) ([]Contact, int, error) {
	limitString := ""
	if limit > 0 {
		limitString = " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit)
	} else {
		limitString = " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(math.MaxUint16)
	}
	queryString := "FOR c in contacts SORT c." + sort + " " + order + limitString + " RETURN c"
	arangoQuery := ara.NewQuery(queryString)
	arangoQuery.SetFullCount(true)
	cursor, err := config.Db().Execute(arangoQuery)

	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	result := make([]Contact, len(cursor.Result))
	err = cursor.FetchBatch(&result)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	return result, cursor.FullCount(), nil
}

// Delete removes the entry for a given email
// Returns true if the contact does not exist int he database
// and an error if the contact could not be deleted
func Delete(email string) (bool, error) {
	contact, err := LoadByEmail(email)
	if contact == nil {
		return true, fmt.Errorf("No contact found for email %v", email)
	}
	if err != nil {
		return false, fmt.Errorf("Error while looking for contact %v: %v", email, err.Error())
	}
	errorMap := config.Context().Delete(contact)
	if value, ok := errorMap["error"]; ok {
		return false, fmt.Errorf("Impossible to delete contact with email %q. %v", email, value)
	}
	return false, nil
}
