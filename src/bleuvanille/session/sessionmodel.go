package session

import (
	"bleuvanille/auth"
	"errors"
	"log"
	"time"

	"github.com/twinj/uuid"
)

// Session stores values for an active user. The user maybe authenticated or not
// How to handle secure sessions: https://www.owasp.org/index.php/Session_Management_Cheat_Sheet
type Session struct {
	ID        uuid.UUID
	UserID    string
	AuthToken string
	IsAdmin   bool
	values    map[string]interface{}
	ExpiresAt time.Time
}

// The duration of the session, in hours
// TODO: set in config file?
const SessionDuration = 1

// NewSession creates a Session instance
func NewSession(userID string, authToken string, isAdmin bool) (Session, error) {
	var session Session

	if userID == "" {
		errorMessage := "Cannot create session, user is missing"
		log.Println(errorMessage)
		return session, errors.New(errorMessage)
	}

	uniqueID := uuid.NewV4()

	stringValueStore := make(map[string]interface{})

	expirationDate := time.Now().Add(time.Hour * SessionDuration)
	session = Session{uniqueID, userID, authToken, isAdmin, stringValueStore, expirationDate}

	return session, nil
}

//IsAuthenticated returns true if the user is connected
func (session *Session) IsAuthenticated() bool {
	if session.AuthToken == "" {
		return false
	}
	token, err := auth.ExtractToken(session.AuthToken)
	if err != nil || token == nil {
		return false
	}

	return token.Valid
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
