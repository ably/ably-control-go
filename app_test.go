package control

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	apps, err := client.Apps()
	assert.NoError(t, err)

	assert.NotEqual(t, len(apps), 0)

	a, err := client.UpdateApp(app.ID, &App{TLSOnly: false})
	assert.NoError(t, err)
	assert.False(t, a.TLSOnly)
	newApp := App{
		Status:                 "disabled",
		TLSOnly:                true,
		ApnsUseSandboxEndpoint: true,
	}
	a, err = client.UpdateApp(app.ID, &newApp)
	assert.NoError(t, err)

	assert.Equal(t, newApp.Status, a.Status)
	assert.Equal(t, newApp.TLSOnly, a.TLSOnly)
	assert.Equal(t, newApp.ApnsUseSandboxEndpoint, a.ApnsUseSandboxEndpoint)

	err = client.DeleteApp(app.ID)
	assert.NoError(t, err)
}
