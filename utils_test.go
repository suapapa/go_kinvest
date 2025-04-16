package kinvest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAccout(t *testing.T) {
	account := "12345678-12"
	first, second, err := parseAccount(account)
	assert.NoError(t, err)

	assert.NotNil(t, first)
	assert.NotNil(t, second)
	assert.Equal(t, "12345678", *first, "first part should be 12345678")
	assert.Equal(t, 12, *second, "second part should be 12")
}
