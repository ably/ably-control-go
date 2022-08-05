package control

// Region is an enum of the possible queue regions.
type Region string

// UsEast1A is the us east 1 a region.
const UsEast1A Region = "us-east-1-a"

// EuWest1A is the eu west 1 a region.
const EuWest1A Region = "eu-west-1-a"

// Amqp contains a queue's amqp data.
type Amqp struct {
	// URI for the AMQP queue interface.
	Uri string `json:"uri,omitempty"`
	// Name of the Ably queue.
	QueueName string `json:"queueName,omitempty"`
}

// Stomp contains a queue's stomp data.
type Stomp struct {
	// URI for the STOMP queue interface.
	Uri string `json:"uri,omitempty"`
	// The host type for the queue.
	Host string `json:"host,omitempty"`
	// Destination queue.
	Destination string `destination:"uri,omitempty"`
}

// Messages contains messages in a queue.
type Messages struct {
	// The number of ready messages in the queue.
	Ready int `json:"ready"`
	// The number of unacknowledged messages in the queue.
	Unacknowledged int `json:"unacknowledged"`
	// The total number of messages in the queue..
	Total int `json:"total"`
}

// Stats contains statistics about an Ably queue
type Stats struct {
	// The rate at which messages are published to the queue. Rate is messages per minute.
	PublishRate int `json:"publishRate"`
	// The rate at which messages are delivered from the queue. Rate is messages per minute.
	DeliveryRate int `json:"deliveryRate"`
	// The rate at which messages are acknowledged. Rate is messages per minute.
	AcknowledgementRate int `json:"acknowledgementRate"`
}

// Queue represents an Ably queue.
type Queue struct {
	// The ID of the Ably queue.
	ID string `json:"id,omitempty"`
	// The Ably application ID.
	AppID string `json:"appId,omitempty"`
	// The friendly name of the queue.
	Name string `json:"name,omitempty"`
	// The data center region for the queue.
	Region Region `json:"region,omitempty"`
	// The amqp data.
	Amqp Amqp `json:"amqp"`
	// The stomp data.
	Stomp Stomp `json:"stomp"`
	// The current state of the queue.
	State string `json:"state,omitempty"`
	// Details of messages in the queue.
	Messages Messages `json:"messages"`
	// Queue stats.
	Stats Stats `json:"stats"`
	// TTL in minutes.
	Ttl int `json:"ttl"`
	// Message limit in number of messages.
	MaxLength int `json:"maxLength"`
	// A boolean that indicates whether this is a dead letter queue or not.
	DeadLetter bool `json:"deadLetter"`
	// The ID of the dead letter queue.
	DeadLetterID string `json:"deadLetterId,omitempty"`
}

// NewQueue is used to create a new Ably queue.
type NewQueue struct {
	// The friendly name of the queue.
	Name string `json:"name,omitempty"`
	// TTL in minutes.
	Ttl int `json:"ttl"`
	// Message limit in number of messages.
	MaxLength int `json:"maxLength"`
	// The data center region for the queue.
	Region Region `json:"region,omitempty"`
}

// Queues lists the queues associated with the specified application ID.
func (c *Client) Queues(appID string) ([]Queue, error) {
	var queues []Queue
	err := c.request("GET", "/apps/"+appID+"/queues", nil, &queues)
	return queues, err
}

// CreateQueue creates a queue for the application specified by application ID.
func (c *Client) CreateQueue(appID string, queue *NewQueue) (Queue, error) {
	var out Queue
	err := c.request("POST", "/apps/"+appID+"/queues", &queue, &out)
	return out, err
}

// DeleteQueue delete the queue with the specified queue name, from the application with the specified application ID.
func (c *Client) DeleteQueue(appID, queueID string) error {
	err := c.request("DELETE", "/apps/"+appID+"/queues/"+queueID, nil, nil)
	return err
}
