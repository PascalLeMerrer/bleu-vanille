package contact

import (
	"bleuvanille/config"
	"errors"
	"log"
	"time"

	ara "github.com/diegogub/aranGO"
)

// Contact is a user who registered to get information about our service
type Contact struct {
	ara.Document
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UserAgent string
	Referer   string
	TimeSpent int
}

// Contacts is a list of Contact
type Contacts []Contact

// New creates a Contact instance
func New(email string, userAgent string, referer string, timeSpent int) (Contact, error) {
	var contact Contact
	if email == "" {
		errorMessage := "Cannot create contact, email is missing"
		log.Println(errorMessage)
		return contact, errors.New(errorMessage)
	}
	contact.Email = email
	contact.CreatedAt = time.Now()
	contact.UserAgent = userAgent
	contact.Referer = referer
	contact.TimeSpent = timeSpent

	return contact, nil
}

// GetKey returns the key in ArangoDB for the contact
func (contact *Contact) GetKey() string {
	return contact.Key
}

// GetCollection returns the collection name in ArangoDB for contacts
func (contact *Contact) GetCollection() string {
	return config.COLNAME_CONTACTS
}

// GetError returns true if there is an error and gives the last error message
func (contact *Contact) GetError() (string, bool) {
	return contact.Message, contact.Error
}
