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

func (e *Contact) GetKey() string{
  return e.Key
}

func (e *Contact) GetCollection() string {
  return config.COLNAME_CONTACTS
}

func (e *Contact) GetError()(string,bool){
    // default error bool and messages. Could be any kind of error
    return e.Message,e.Error
}