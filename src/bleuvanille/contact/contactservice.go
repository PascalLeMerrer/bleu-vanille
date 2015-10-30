package contact

// Ensures persistance of contacts in Postgresql database

import (
	"bleuvanille/config"
	"errors"
	"fmt"

	ara "github.com/diegogub/aranGO"
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
func LoadAll() ([]Contact, error) {
	queryString := "FOR c in contacts RETURN c"

	arangoQuery := ara.NewQuery(queryString)
	cursor, err := config.Db().Execute(arangoQuery)

	if err != nil {
		return nil, err
	}
	result := make([]Contact, len(cursor.Result))
	err = cursor.FetchBatch(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete removes the entry for a given email
func Delete(email string) error {
	contact, err := LoadByEmail(email)
	if contact == nil {
		return fmt.Errorf("No contact found for email %v", email)
	}
	if err != nil {
		return fmt.Errorf("Error while looking for contact %v: %v", email, err.Error())
	}
	errorMap := config.Context().Delete(contact)
	if value, ok := errorMap["error"]; ok {
		return fmt.Errorf("Impossible to delete contact by email %q because %v", email, value)
	}
	return nil
}
