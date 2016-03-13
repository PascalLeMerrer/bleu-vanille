package session

import (
	"bleuvanille/config"
	"errors"
	"log"
	"time"
)

// Session stores values for an active user. The user maybe authenticated or not
// How to handle secure sessions: https://www.owasp.org/index.php/Session_Management_Cheat_Sheet
type Session struct {
	ID        string `json:"id"`
	Key       string `json:"_key,omitempty"`
	SessionID string
	UserID    string
	Email     string
	Firstname string
	Lastname  string
	IsAdmin   bool
	values    map[string]interface{}
	ExpiresAt time.Time
}

// New creates a Session instance
func New(sessionID string, userID string, firstname string, lastname string, email string, isAdmin bool) (Session, error) {
	var session Session

	if userID == "" {
		errorMessage := "Error: Cannot create session, user is missing"
		log.Println(errorMessage)
		return session, errors.New(errorMessage)
	}

	stringValueStore := make(map[string]interface{})

	expirationDate := time.Now().Add(time.Hour * config.SessionDuration)
	session = Session{SessionID: sessionID, UserID: userID, Email: email, Firstname: firstname, Lastname: lastname, IsAdmin: isAdmin, values: stringValueStore, ExpiresAt: expirationDate}
	return session, nil
}

// Set stores a value in the session
// TODO should we really have this? Flash (transient) values may be better
func (session *Session) Set(key string, value interface{}) {
	session.values[key] = value
}

// Get returns a value stored in the session or nil if this value does not exist
func (session *Session) Get(key string) interface{} {
	return session.values[key]
}

// GetKey returns the primary key
func (session *Session) GetKey() string {
	return session.Key
}
