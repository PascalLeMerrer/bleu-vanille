package contact

// Ensures persistance of contacts in Postgresql database

import (
	"bleuvanille/config"
	"log"
)

// Save inserts a contact into the database
func Save(contact Contact) (Contact, error) {
	_, err := config.Db().Query("INSERT INTO contacts (email, created_at) VALUES ($1,$2);", contact.Email, contact.CreatedAt)
	return contact, err
}

// LoadByEmail returns the contact object for a given email, if any
func LoadByEmail(email string) (*Contact, error) {
	var contact Contact
	row := config.Db().QueryRow("SELECT * FROM contacts WHERE email = $1;", email)
	err := row.Scan(&contact.Email, &contact.CreatedAt)
	if err != nil {
		log.Println(err)
	}
	return &contact, err
}

// LoadAll returns the list of all contacts in the database
func LoadAll() (*Contacts, error) {
	var contacts Contacts

	rows, err := config.Db().Query("SELECT * FROM contacts")

	if err != nil {
		log.Printf("Cannot query contact list: %s", err)
		return nil, err
	}

	// Close rows after all readed
	defer rows.Close()

	for rows.Next() {
		var c Contact

		err := rows.Scan(&c.Email, &c.CreatedAt)

		if err != nil {
			log.Println(err)
		}

		contacts = append(contacts, c)
	}

	return &contacts, err
}

// Delete deletes the entry for a given email
func Delete(email string) error {
	stmt, err := config.Db().Prepare("DELETE FROM contacts WHERE email=$1;")

	if err != nil {
		log.Printf("Cannot delete contact: %s", err)
		return err
	}

	_, err = stmt.Exec(email)

	return err
}
