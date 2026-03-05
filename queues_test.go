package control

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueues(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	name := "queue-key-" + fmt.Sprint(rand.Uint64())

	queue := NewQueue{
		Name:      name,
		Ttl:       50,
		MaxLength: 10,
		Region:    EuWest1A,
	}

	q, err := client.CreateQueue(app.ID, &queue)
	require.NoError(t, err)
	assert.Equal(t, queue.Name, q.Name)
	assert.Equal(t, queue.Ttl, q.Ttl)
	assert.Equal(t, queue.MaxLength, q.MaxLength)
	assert.Equal(t, queue.Region, q.Region)

	queues, err := client.Queues(app.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, queues)

	err = client.DeleteQueue(app.ID, q.ID)
	assert.NoError(t, err)
}

func TestStompDestinationJsonRoundTrip(t *testing.T) {
	original := Stomp{
		Uri:         "stomp://localhost:61613",
		Host:        "my-host",
		Destination: "/queue/test",
	}

	data, err := json.Marshal(original)
	assert.NoError(t, err)

	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	assert.NoError(t, err)
	assert.Contains(t, raw, "destination", "Stomp.Destination should marshal with json key 'destination'")
	assert.Equal(t, "/queue/test", raw["destination"])

	var decoded Stomp
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}
