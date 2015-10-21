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
	//Verify if the email already exists
	existingContact, err := LoadByEmail(contact.Email)

	if err != nil {
		strerr := fmt.Sprintf("%q", err)
		return errors.New(strerr)
	}

	if existingContact == nil || len(existingContact.Key) == 0 {
		err := config.Context().Save(contact)

		if err != nil {
			strerr := fmt.Sprintf("%q", err)
			return errors.New(strerr)
		}
	}

	return nil
}

// LoadByEmail returns the contact object for a given email, if any
func LoadByEmail(email string) (*Contact, error) {
	var result Contact

	col := config.GetCollection(&result)
	result.Email = email
	cursor, err := col.Example(result, 0, 1)

	//return an error
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
func LoadAll() (*[]Contact, error) {
	var contact Contact
	querystr := fmt.Sprintf("FOR c in contacts RETURN c")

	arangoQuery := ara.NewQuery(querystr)
	cursor, err := config.Db().Execute(arangoQuery)

	//return an error
	if err != nil {
		return nil, err
	}

	if cursor.Result != nil && len(cursor.Result) > 0 {
		result := make([]Contact, len(cursor.Result))

		for cursor.FetchOne(&contact) {
			result = append(result, contact)
		}

		return &result, nil
	}

	return nil, nil
}

// Delete deletes the entry for a given email
func Delete(email string) error {
	contact, _ := LoadByEmail(email)

	if contact == nil {
		return nil
	}

	err := config.Context().Delete(contact)

	if err != nil  && len(err) > 0 {
		errorstring := fmt.Sprintf("Impossible to delete contact by email %q because %s", email, err)
		return errors.New(errorstring)
	}

	return nil
}
