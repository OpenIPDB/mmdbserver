package mmdbserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneralRequest(t *testing.T) {
	request, err := NewRequest(mustFixtureRequest("general"))
	assert.NoError(t, err)
	assert.True(t, request.NewDeploy())
}

func TestInvalidRequest(t *testing.T) {
	fixtures := map[string]error{
		"invalid-edition-id":  ErrInvalidEditionID,
		"invalid-account-id":  ErrInvalidAccount,
		"invalid-license-key": ErrInvalidAccount,
		"invalid-hash":        ErrInvalidMD5Hash,
	}
	for name, target := range fixtures {
		_, err := NewRequest(mustFixtureRequest(name))
		assert.ErrorIs(t, err, target)
	}
}
