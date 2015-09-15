package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGenerationAndExtraction(t *testing.T) {
	encodedToken, sessionID := GetEncodedToken()
	assert.NotEmpty(t, encodedToken, "No encoded token generated.")
	assert.True(t, len(encodedToken) >= 16, "invalid ")

	assert.NotEmpty(t, sessionID, "No sessionID generated.")
	assert.True(t, len(sessionID) >= 16, "invalid ")

	token, err := ExtractToken(encodedToken)
	assert.NoError(t, err, "Error while generated token.")
	assert.True(t, token.Valid, "The generated token is not valid")
}
