package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestEmail = "tester@email.org"
)

func TestSessionToken(t *testing.T) {
	encodedToken, sessionID := GetSessionToken()
	assert.NotEmpty(t, encodedToken, "No encoded token generated.")
	assert.True(t, len(encodedToken) >= 16, "invalid ")

	assert.NotEmpty(t, sessionID, "No sessionID generated.")
	assert.True(t, len(sessionID) >= 16, "invalid ")

	token, err := ExtractToken(encodedToken)
	assert.NoError(t, err, "Error while generated token.")
	assert.True(t, token.Valid, "The generated token is not valid")
	assert.Equal(t, sessionID, token.Claims["id"], "The token should contain the session ID")
}

func TestResetToken(t *testing.T) {
	encodedToken := GetResetToken(TestEmail)
	assert.NotEmpty(t, encodedToken, "No encoded token generated.")
	assert.True(t, len(encodedToken) >= 16, "invalid ")

	token, err := ExtractToken(encodedToken)
	assert.NoError(t, err, "Error while generated token.")
	assert.True(t, token.Valid, "The generated token is not valid")
	assert.Equal(t, TestEmail, token.Claims["email"], "The token should contain the user email")
}
