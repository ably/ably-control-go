package control

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	apps, err := client.Apps()
	require.NoError(t, err)
	assert.NotEmpty(t, apps)

	a, err := client.UpdateApp(app.ID, &NewApp{TLSOnly: false})
	require.NoError(t, err)
	assert.False(t, a.TLSOnly)

	newApp := NewApp{
		Status:                 "disabled",
		TLSOnly:                true,
		ApnsUseSandboxEndpoint: true,
	}
	a, err = client.UpdateApp(app.ID, &newApp)
	require.NoError(t, err)
	assert.Equal(t, newApp.Status, a.Status)
	assert.Equal(t, newApp.TLSOnly, a.TLSOnly)
	assert.Equal(t, newApp.ApnsUseSandboxEndpoint, a.ApnsUseSandboxEndpoint)
}

func TestAppUnmarshalIncludesNewFields(t *testing.T) {
	data := `{
		"id": "app123",
		"accountId": "acc456",
		"name": "Test App",
		"status": "enabled",
		"tlsOnly": true,
		"created": 1700000000,
		"modified": 1700001000,
		"_links": {
			"self": "https://control.ably.net/v1/apps/app123"
		}
	}`

	var app App
	err := json.Unmarshal([]byte(data), &app)
	assert.NoError(t, err)
	assert.Equal(t, "app123", app.ID)
	assert.Equal(t, int64(1700000000), app.Created)
	assert.Equal(t, int64(1700001000), app.Modified)
	assert.NotNil(t, app.Links)
	assert.Equal(t, "https://control.ably.net/v1/apps/app123", app.Links["self"])
}
