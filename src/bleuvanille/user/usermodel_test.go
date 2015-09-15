package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TestEmail = "me@mail.org"
const TestFirstname = "John"
const TestLastname = "Doe"
const TestPassword = "Pa%%W0rD"

func TestUserCreation(t *testing.T) {
	user, userCreationError := New(TestEmail, TestFirstname, TestLastname, TestPassword)

	assert.NoError(t, userCreationError, "User creation error.")
	assert.NotNil(t, user, "User creation error.")

	assert.Equal(t, TestEmail, user.Email, "User email not set during user creation")
	assert.Equal(t, TestFirstname, user.Firstname, "User firstname not set during user creation")
	assert.Equal(t, TestLastname, user.Lastname, "User lastname not set during user creation")
	assert.NotEmpty(t, user.Hash, "User password hash not set during user creation")
	assert.NotEmpty(t, user.ID, "User ID not set during user creation")
	assert.True(t, len(user.Hash) >= 8, "User password hash not long enough")
	assert.True(t, len(user.ID) >= 8, "User ID not long enough")
	assert.False(t, user.IsAdmin, "New user must not be admin")
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Minute, "User creation date not set correctly during user creation")

}

func TestIntToIntArray(t *testing.T) {

	result := intToIntArray(1234567890123456, 8)
	assert.Len(t, result, 8, "intToIntArray should return the number of element given as a second parameter")
	assert.Equal(t, 12, result[0])
	assert.Equal(t, 34, result[1])
	assert.Equal(t, 56, result[2])
	assert.Equal(t, 78, result[3])
	assert.Equal(t, 90, result[4])
	assert.Equal(t, 12, result[5])
	assert.Equal(t, 34, result[6])
	assert.Equal(t, 56, result[7])
}

func TestGenerateID(t *testing.T) {
	ID := generateID()
	assert.True(t, len(ID) > 8, "ID should be at least 8 characters long")
}
