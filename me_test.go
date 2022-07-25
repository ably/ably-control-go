package control

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMe(t *testing.T) {
	_, me := newTestClient(t)

	assert.NotEmpty(t, me.Account.ID)
	assert.NotEmpty(t, me.Account.Name)
	assert.NotEmpty(t, me.Token.ID)
	assert.NotEmpty(t, me.User.ID)
}
