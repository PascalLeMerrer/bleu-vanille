package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGenerationAndExtraction(t *testing.T) {
	encodedToken := GetEncodedToken()
	assert.NotEmpty(t, encodedToken, "No encoded token generated.")
	assert.True(t, len(encodedToken) >= 16, "invalid ")

	token, err := ExtractToken(encodedToken)
	assert.NoError(t, err, "Error while generated token.")
	assert.True(t, token.Valid, "The generated token is not valid")
}
