package contact

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicateEmail(t *testing.T) {
	var email = "contactDuplicateTest@mail.com"
	contact := Contact{Email: email}
	err := Save(&contact)
	assert.NoError(t, err, "Saving contact")
	contact2 := Contact{Email: email}
	err = Save(&contact2)
	assert.Error(t, err, "Saving contact with duplicate email address failed")
	err = Remove(email)
	assert.Nil(t, err, "Deleting non existant email should not return an error")
}

func TestDeleteUnknownEmail(t *testing.T) {
	var email = "contactDeleteTest@unknown.com"
	err := Remove(email)
	assert.Nil(t, err, "Deleting non existant email should not return an error")
}

func TestSearchByEmail(t *testing.T) {
	var email = "contactSearchTest@bleuvanille.com"
	created := Contact{Email: email}

	//Delete if any
	Remove(email)

	//Save the new contact
	Save(&created)

	//Search it afterwards
	c, err := FindByEmail(email)

	assert.NoError(t, err, "Load By Email error.")
	assert.NotEqual(t, c, nil, "User is not found by email")

	assert.Equal(t, email, c.Email, "Wrong email.")

	//purge the database with the test contact
	Remove(email)
}

func TestGetAll(t *testing.T) {
	contacts, count, err := FindAll("email", "ASC", 0, 0)
	assert.NoError(t, err, "Error when loading contacts")

	initialCount := len(contacts)

	var email1 = "testGetAll1@bleuvanille.com"
	var email2 = "testGetAll2@bleuvanille.com"
	contact1 := Contact{Email: email1}
	contact2 := Contact{Email: email2}

	Save(&contact1)
	Save(&contact2)

	contacts, count, err = FindAll("email", "ASC", 0, 0)
	assert.Equal(t, len(contacts), count, "Inconsistant data returned\n")
	assert.Equal(t, count, initialCount+2, "Contacts not added in %v\n", contacts)

	email1Found := false
	email2Found := false
	for _, contact := range contacts {
		if contact.Email == email1 {
			if email1Found {
				assert.Fail(t, "Duplicate contact %v\n", email1)
			}
			email1Found = true
		}
		if contact.Email == email2 {
			if email2Found {
				assert.Fail(t, "Duplicate contact %v\n", email2)
			}
			email2Found = true
		}
	}
	assert.True(t, email1Found, "Wrong email for contact %v\n", email1)
	assert.True(t, email2Found, "Wrong email for contact %v\n", email2)

	assert.NoError(t, err, "Error when loading contacts")

	Remove(email1)
	Remove(email2)

	contacts, count, err = FindAll("email", "ASC", 0, 0)

	assert.Equal(t, len(contacts), initialCount, "Contacts not deleted from %v\n", contacts)
	assert.NoError(t, err, "Contact deletion failed.")

}
