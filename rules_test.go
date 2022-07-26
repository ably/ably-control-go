package control

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRules(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)
	defer client.DeleteApp(app.ID)

	target := &HttpTarget{
		Url:       "http://test.com",
		Headers:   []Header{{Name: "a", Value: "b"}},
		Enveloped: true,
		Format:    "json",
	}

	rule := NewRule{
		Status:      "enabled",
		RequestMode: "single",
		Source: Source{
			ChannelFilter: "aaa",
			Type:          "channel.message",
		},
		Target: target,
	}

	r, err := client.CreateRule(app.ID, &rule)
	assert.NoError(t, err)
	assert.Equal(t, rule.Status, r.Status)

	r2, err := client.Rule(app.ID, r.ID)
	assert.NoError(t, err)
	assert.Equal(t, r, r2)

	err = client.DeleteRule(app.ID, r.ID)
	assert.NoError(t, err)

	err = client.DeleteApp(app.ID)
	assert.NoError(t, err)
}
