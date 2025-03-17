package control

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, namespace, n)

	namespaces, err := client.Namespaces(app.ID)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.Equal(t, namespace, n)

	err = client.DeleteNamespace(app.ID, namespace.ID)
	assert.NoError(t, err)

	err = client.DeleteApp(app.ID)
	assert.NoError(t, err)
}
