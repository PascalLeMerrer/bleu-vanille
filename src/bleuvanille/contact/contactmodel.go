package contact

import (
	"errors"
	"log"
	"time"
)

// Contact is a user who registered to get information about our service
type Contact struct {
	ID        string    `json:"id"`
	Key       string    `json:"_key,omitempty"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UserAgent string    `json:"useragent,omitempty"`
	Referer   string
	TimeSpent int `json:"timespent,omitempty"`
}

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
