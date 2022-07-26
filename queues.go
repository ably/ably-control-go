package control

type Amqp struct {
	Uri       string `json:"uri,omitempty"`
	QueueName string `json:"queueName,omitempty"`
}

type Stomp struct {
	Uri         string `json:"uri,omitempty"`
	Host        string `json:"host,omitempty"`
	Destination string `destination:"uri,omitempty"`
}

type Messages struct {
	Ready          int `json:"ready"`
	Unacknowledged int `json:"unacknowledged"`
	Total          int `json:"total"`
}

type Stats struct {
	PublishRate         int `json:"publishRate"`
	DeliveryRate        int `json:"deliveryRate"`
	AcknowledgementRate int `json:"acknowledgementRate"`
}

type Queue struct {
	ID           string   `json:"id,omitempty"`
	AppID        string   `json:"appId,omitempty"`
	Name         string   `json:"name,omitempty"`
	Region       string   `json:"region,omitempty"`
	Amqp         Amqp     `json:"amqp"`
	Stomp        Stomp    `json:"stomp"`
	State        string   `json:"state,omitempty"`
	Messages     Messages `json:"messages"`
	Stats        Stats    `json:"stats"`
	Ttl          int      `json:"ttl"`
	MaxLength    int      `json:"maxLength"`
	DeadLetter   bool     `json:"deadLetter"`
	DeadLetterID string   `json:"deadLetterId,omitempty"`
}

type NewQueue struct {
	Name      string `json:"name,omitempty"`
	Ttl       int    `json:"ttl"`
	MaxLength int    `json:"maxLength"`
	Region    string `json:"region,omitempty"`
}

func (c *Client) Queues(appID string) ([]Queue, error) {
	var queues []Queue
	err := c.request("GET", "/apps/"+appID+"/queues", nil, &queues)
	return queues, err
}

func (c *Client) CreateQueue(appID string, queue *NewQueue) (Queue, error) {
	var out Queue
	err := c.request("POST", "/apps/"+appID+"/queues", &queue, &out)
	return out, err
}

func (c *Client) DeleteQueue(appID, queueID string) error {
	err := c.request("DELETE", "/apps/"+appID+"/queues/"+queueID, nil, nil)
	return err
}
