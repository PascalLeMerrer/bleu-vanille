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

	// Cleanup before testing
	_ = Delete(&testUser)

	userSaveError := Save(&testUser)
	assert.NoError(t, userSaveError, "User save error.")

	testSession, sessionCreationError := session.New(testSessionID, testUser.ID, testUser.IsAdmin)
	assert.NoError(t, sessionCreationError, "Session creation error.")
	err := session.Save(&testSession)
	assert.NoError(t, err, "Session save error.")

	userDeletionError := Delete(&testUser)
	assert.NoError(t, userDeletionError, "User Deletion error.")

}


func TestUserUpdate(t *testing.T) {
	testUser, userCreationError := New(TestEmail, TestFirstname, TestLastname, TestPassword)
	assert.NoError(t, userCreationError, "User creation error.")

	// Cleanup before testing
	_ = Delete(&testUser)

	userSaveError := Save(&testUser)
	assert.NoError(t, userSaveError, "User save error.")

	testUser.Email = "newemail@bleuvanille.com"
	testUser.IsAdmin = true
	testUser.Firstname = "Jane"
	testUser.Lastname = "Groquik"
	testUser.Hash = "NewHash"
	testUser.ResetToken = "a new token"

	userUpdateError := Update(&testUser)
	assert.NoError(t, userUpdateError, "User update error.")

	updatedUser, loadUserError := LoadByID(testUser.ID)
	assert.NoError(t, loadUserError, "User loading error.")
	assert.Equal(t, testUser.Email, updatedUser.Email, "Updated user does not contain the expected Email.")
	assert.Equal(t, testUser.Firstname, updatedUser.Firstname, "Updated user does not contain the expected Firstname.")
	assert.Equal(t, testUser.Hash, updatedUser.Hash, "Updated user does not contain the expected Hash.")
	assert.Equal(t, testUser.IsAdmin, updatedUser.IsAdmin, "Updated user does not contain the expected IsAdmin.")
	assert.Equal(t, testUser.Lastname, updatedUser.Lastname, "Updated user does not contain the expected Lastname.")
	assert.Equal(t, testUser.ID, updatedUser.ID, "Updated user does not contain the expected ID.")
	assert.Equal(t, testUser.ResetToken, updatedUser.ResetToken, "Updated user does not contain the expected ResetToken.")

	// Cleanup after testing
	userDeletionError := Delete(&testUser)
	assert.NoError(t, userDeletionError, "User Deletion error.")
}
