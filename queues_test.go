package control

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueues(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	name := "queue-key-" + fmt.Sprint(rand.Uint64())

	queue := NewQueue{
		Name:      name,
		Ttl:       50,
		MaxLength: 10,
		Region:    "eu-west-1-a",
	}

	q, err := client.CreateQueue(app.ID, &queue)
	assert.NoError(t, err)
	assert.Equal(t, queue.Name, q.Name)
	assert.Equal(t, queue.Ttl, q.Ttl)
	assert.Equal(t, queue.MaxLength, q.MaxLength)
	assert.Equal(t, queue.Region, q.Region)

	queues, err := client.Queues(app.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, queues)

	queue = NewQueue{
		Name:      name + "-changed",
		Ttl:       40,
		MaxLength: 20,
		Region:    "us-east-1-a",
	}

	err = client.DeleteQueue(app.ID, q.ID)
	assert.NoError(t, err)

	err = client.DeleteApp(app.ID)
	assert.NoError(t, err)
}
