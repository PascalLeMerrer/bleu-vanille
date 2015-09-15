package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const key1 = "key1"
const key2 = "key2"
const key3 = "key3"
const value1 = "a string"
const value2 = true
const value3 = 123
const testUserID = "ABCD123"
const testSessionID = "dlùaskda¨$ùùllk1"
const testIsAdmin = false

func TestSessionCreation(t *testing.T) {

	session, sessionCreationError := New(testSessionID, testUserID, testIsAdmin)
	assert.NoError(t, sessionCreationError, "Session creation error.")
	assert.Equal(t, testSessionID, session.SessionID, "SessionID not set in session.")
	assert.Equal(t, testUserID, session.UserID, "User not set in session.")
	assert.Equal(t, testUserID, session.UserID, "User not set in session.")
	assert.False(t, session.IsAdmin, "User should not declared as admin in session.")
	assert.WithinDuration(t, time.Now().Add(time.Hour), session.ExpiresAt, time.Minute, "Invalid session duration.")
}

func TestSessionValueStore(t *testing.T) {
	session, sessionCreationError := New(testSessionID, testUserID, testIsAdmin)
	assert.NoError(t, sessionCreationError, "Session creation error.")
	assert.NotNil(t, session, "Session creation error.")

	session.Set(key1, value1)
	session.Set(key2, value2)
	session.Set(key3, value3)
	assert.Equal(t, value1, session.Get(key1), "Session should return a string.")
	assert.Equal(t, value2, session.Get(key2), "Session should return a boolean.")
	assert.Equal(t, value3, session.Get(key3), "Session should return an integer.")
}
