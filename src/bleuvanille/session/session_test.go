package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchByID(t *testing.T) {
	session, _ := FindByID("UNKNOWNID")

	if session != nil {
		t.Errorf("Found a session that does not exist : %+v", session)
	}

	var session1 Session
	session1.SessionID = "1"
	session1.UserID = "1"
	session1.IsAdmin = false

	Save(&session1)

	sessionbyid, err := FindByID("1")

	assert.NoError(t, err, "Get Session by ID error.")

	if sessionbyid == nil {
		t.Errorf("Not Found a session with the id \"1\"")
	} else {
		if sessionbyid.SessionID != "1" {
			t.Errorf("Found a session with the wrong id : %+v", session)
		}
	}

	//purge the database with the test session
	Remove("1")

	session, _ = FindByID("1")

	if session != nil {
		t.Errorf("Session incorrectly deleted")
	}

}
