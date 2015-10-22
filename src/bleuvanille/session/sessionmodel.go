package session

import (
	"bleuvanille/config"
	"errors"
	"log"
	"time"
	ara "github.com/diegogub/aranGO"
)

// Session stores values for an active user. The user maybe authenticated or not
// How to handle secure sessions: https://www.owasp.org/index.php/Session_Management_Cheat_Sheet
type Session struct {
	ara.Document
	SessionID string
	UserID    string
	IsAdmin   bool
	values    map[string]interface{}
	ExpiresAt time.Time
}

// New creates a Session instance
func New(sessionID string, userID string, isAdmin bool) (Session, error) {
	var session Session

	if userID == "" {
		errorMessage := "Error: Cannot create session, user is missing"
		log.Println(errorMessage)
		return session, errors.New(errorMessage)
	}

	stringValueStore := make(map[string]interface{})

	expirationDate := time.Now().Add(time.Hour * config.SessionDuration)
	session = Session{SessionID : sessionID, UserID : userID, IsAdmin : isAdmin, values:stringValueStore, ExpiresAt : expirationDate}

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

func (e *Session) GetKey() string{
  return e.Key
}

func (e *Session) GetCollection() string {
  return config.COLNAME_SESSIONS
}

func (e *Session) GetError()(string,bool){
    // default error bool and messages. Could be any kind of error
    return e.Message,e.Error
}