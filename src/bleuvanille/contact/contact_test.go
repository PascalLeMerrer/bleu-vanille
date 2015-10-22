package contact

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSearchEmail(t *testing.T) {
	var email = "test@testdecreation.com"
	ccreated  := Contact{ Email : email}

	//Delete if any
	Delete(email)
	
	//Save the new contact
	Save(&ccreated)
	
	//Search it afterwards
	c, err := LoadByEmail(email)
	
	assert.NoError(t, err, "Load By Email error.")
	assert.NotEqual(t, c, nil, "User is not found by email")
	
	if c.Email != email {
		t.Errorf("Not Found correct email : %q", c.Email)
	}

	//purge the database with the test contact
	Delete(email)
}

func TestGetAll(t *testing.T) {
	var email1 = "test1@testdecreation.com"
	var email2 = "test2@testdecreation.com"
	contactCreated1  := Contact{ Email : email1}
	contactCreated2  := Contact{ Email : email2}
	//Save the new contact
	Save(&contactCreated1)
	Save(&contactCreated2)
	
	//Search it afterwards
	contacts, _ := LoadAll()
	
	if len(*contacts) < 1 {
		t.Errorf("Found no contact", contacts)
	}
	
	Delete(email1)
	Delete(email2)
}
