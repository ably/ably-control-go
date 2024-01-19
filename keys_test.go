package control

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	name := "test-key-" + fmt.Sprint(rand.Uint64())

	key := NewKey{
		Name:            name,
		Capability:      map[string][]string{"a": {"subscribe"}},
		RevocableTokens: true,
	}

	k, err := client.CreateKey(app.ID, &key)
	assert.NoError(t, err)
	assert.Equal(t, key.Name, k.Name)
	assert.Equal(t, key.Capability, k.Capability)
	assert.Equal(t, k.Status, 0)
	assert.Equal(t, key.RevocableTokens, k.RevocableTokens)
	assert.NotEmpty(t, k.AppID)
	assert.NotEmpty(t, k.Created)
	assert.NotEmpty(t, k.Modified)
	assert.NotEmpty(t, k.ID)
	assert.NotEmpty(t, k.Key)

	keys, err := client.Keys(app.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, keys)

	key = NewKey{
		Name:            name + "-changed",
		Capability:      map[string][]string{"b": {"publish"}},
		RevocableTokens: false,
	}

	k, err = client.UpdateKey(app.ID, k.ID, &key)
	assert.NoError(t, err)
	assert.Equal(t, key.Name, k.Name)
	assert.Equal(t, key.Capability, k.Capability)
	assert.Equal(t, key.RevocableTokens, k.RevocableTokens)

	err = client.RevokeKey(app.ID, k.ID)
	assert.NoError(t, err)

	err = client.DeleteApp(app.ID)
	assert.NoError(t, err)
}
