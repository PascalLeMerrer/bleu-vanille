package contact

import (
	"errors"
	"log"
	"time"
)

// Contact is a user who registered to get information about our service
type Contact struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Contacts is a list of Contact
type Contacts []Contact

// New creates a Contact instance
func New(email string) (Contact, error) {
	var contact Contact
	if email == "" {
		errorMessage := "Cannot create contact, email is missing"
		log.Println(errorMessage)
		return contact, errors.New(errorMessage)
	}
	contact.Email = email
	contact.CreatedAt = time.Now()

	return contact, nil
}
