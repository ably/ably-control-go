package control

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var token string
var apps []string

func TestMain(m *testing.M) {
	token = os.Getenv("ABLY_ACCOUNT_TOKEN")
	rand.Seed(time.Now().UnixNano())

	if token == "" {
		panic("ABLY_ACCOUNT_TOKEN not set")
	}

	// Attempt to clean up apps if anything went wrong
	client, _, err := NewClient(token)
	if err == nil {
		for _, v := range apps {
			_ = client.DeleteApp(v)
		}
	}

	code := m.Run()
	os.Exit(code)
}

func newTestApp(t *testing.T, client *Client) App {
	n := rand.Uint64()
	name := "test-" + fmt.Sprint(n)
	t.Logf("crating app with name: %s", name)
	app := App{
		Name:   name,
		Status: "enabled",
		//TLSOnly:                false,
		FcmKey:                 "",
		ApnsCertificate:        "",
		ApnsPrivateKey:         "",
		ApnsUseSandboxEndpoint: false,
	}
	app, err := client.CreateApp(&app)

	assert.NoError(t, err)
	apps = append(apps, app.ID)

	return app
}

func newTestClient(t *testing.T) (Client, Me) {
	client, me, err := NewClient(token)
	assert.NoError(t, err)
	return client, me

}
