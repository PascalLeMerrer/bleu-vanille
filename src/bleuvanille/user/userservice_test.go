package user

import (
	"bleuvanille/session"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testSessionID = "123456"
)

func TestSessionSave(t *testing.T) {
	testUser, userCreationError := New(TestEmail, TestFirstname, TestLastname, TestPassword)
	assert.NoError(t, userCreationError, "User creation error.")

	// Clenup before testing
	_ = Delete(testUser)

	userSaveError := Save(testUser)
	assert.NoError(t, userSaveError, "User save error.")

	testSession, sessionCreationError := session.New(testSessionID, testUser.ID, testUser.IsAdmin)
	assert.NoError(t, sessionCreationError, "Session creation error.")
	err := session.Save(testSession)
	assert.NoError(t, err, "Session save error.")

	userDeletionError := Delete(testUser)
	assert.NoError(t, userDeletionError, "User Deletion error.")

}
