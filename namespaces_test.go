package control

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamespaces(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	name := "test-namespace-" + fmt.Sprint(rand.Uint64())

	namespace := Namespace{
		ID:               name,
		Authenticated:    false,
		Persisted:        false,
		PersistLast:      false,
		PushEnabled:      false,
		TlsOnly:          false,
		ExposeTimeserial: false,
		BatchingEnabled:  false,
	}

	n, err := client.CreateNamespace(app.ID, &namespace)
	require.NoError(t, err)
	namespace.AppID = n.AppID
	namespace.Created = n.Created
	namespace.Modified = n.Modified
	assert.Equal(t, namespace, n)

	namespaces, err := client.Namespaces(app.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, namespaces)

	namespace = Namespace{
		ID:               namespace.ID,
		Authenticated:    true,
		Persisted:        true,
		PersistLast:      true,
		PushEnabled:      true,
		TlsOnly:          true,
		ExposeTimeserial: true,
		BatchingEnabled:  true,
		BatchingInterval: Interval(100),
	}

	n, err = client.UpdateNamespace(app.ID, &namespace)
	require.NoError(t, err)
	namespace.AppID = n.AppID
	namespace.Created = n.Created
	namespace.Modified = n.Modified
	assert.Equal(t, namespace, n)

	namespace = Namespace{
		ID:               namespace.ID,
		Authenticated:    true,
		Persisted:        true,
		PersistLast:      true,
		PushEnabled:      true,
		TlsOnly:          true,
		ExposeTimeserial: true,
		BatchingEnabled:  false,
	}

	n, err = client.UpdateNamespace(app.ID, &namespace)
	require.NoError(t, err)
	namespace.AppID = n.AppID
	namespace.Created = n.Created
	namespace.Modified = n.Modified
	assert.Equal(t, namespace, n)

	namespace = Namespace{
		ID:                 namespace.ID,
		Authenticated:      true,
		Persisted:          true,
		PersistLast:        true,
		PushEnabled:        true,
		TlsOnly:            true,
		ExposeTimeserial:   true,
		ConflationEnabled:  true,
		ConflationInterval: Interval(1000),
		ConflationKey:      "test",
	}

	n, err = client.UpdateNamespace(app.ID, &namespace)
	require.NoError(t, err)
	assert.Equal(t, namespace.ConflationEnabled, n.ConflationEnabled)
	assert.Equal(t, namespace.ConflationInterval, n.ConflationInterval)
	assert.Equal(t, namespace.ConflationKey, n.ConflationKey)

	err = client.DeleteNamespace(app.ID, namespace.ID)
	assert.NoError(t, err)
}

func TestNamespaceRoundTrip(t *testing.T) {
	ns := Namespace{
		ID:                      "test-ns",
		Authenticated:           true,
		MutableMessages:         true,
		PopulateChannelRegistry: true,
		AppID:                   "app123",
		Created:                 1700000000,
		Modified:                1700001000,
	}

	data, err := json.Marshal(ns)
	assert.NoError(t, err)

	var decoded Namespace
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, true, decoded.MutableMessages)
	assert.Equal(t, true, decoded.PopulateChannelRegistry)
	assert.Equal(t, "app123", decoded.AppID)
	assert.Equal(t, int64(1700000000), decoded.Created)
	assert.Equal(t, int64(1700001000), decoded.Modified)
}

func TestNamespaceOmitsZeroResponseFields(t *testing.T) {
	ns := Namespace{
		ID:              "test-ns",
		MutableMessages: true,
	}

	data, err := json.Marshal(ns)
	assert.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	assert.NoError(t, err)

	_, hasAppID := raw["appId"]
	assert.False(t, hasAppID, "appId should be omitted when empty")
	_, hasCreated := raw["created"]
	assert.False(t, hasCreated, "created should be omitted when zero")
	_, hasModified := raw["modified"]
	assert.False(t, hasModified, "modified should be omitted when zero")
}
