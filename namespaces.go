package control

// A struct representing an Ably namespace.
type Namespace struct {
	//The namespace or channel name that the channel rule will apply to. For example,
	//if you specify namespace the namespace will be set to namespace and will match
	//with channels namespace:* and namespace.
	ID string `json:"id,omitempty"`
	// If true, clients will not be permitted to use (including to attach, publish, or subscribe)
	// any channels within this namespace unless they are identified, that is, authenticated using
	// a client ID. See the Ably knowledge base for more details. https://knowledge.ably.com/authenticated-and-identified-clients
	Authenticated bool `json:"authenticated"`
	// If true, all messages on a channel will be stored for 24 hours. You can access stored
	// messages via the History API. Please note that for each message stored, an additional
	// message is deducted from your monthly allocation.
	Persisted bool `json:"persisted"`
	// If true, the last message published on a channel will be stored for 365 days. You can
	// access the stored message only by using the channel rewind mechanism and attaching with rewind=1.
	// Please note that for each message stored, an additional message is deducted from your monthly allocation.
	PersistLast bool `json:"persistLast"`
	// If true, publishing messages with a push payload in the extras field is permitted
	// and can trigger the delivery of a native push notification to registered devices for the channel.
	PushEnabled bool `json:"pushEnabled"`
	// If true, only clients that are connected using TLS will be permitted to subscribe to any
	// channels within this namespace.
	TlsOnly bool `json:"tlsOnly"`
	// If true, messages received on a channel will contain a unique timeserial that can be
	// referenced by later messages for use with message interactions.
	ExposeTimeserial bool `json:"exposeTimeserial"`
	// If true, channels within this namespace will start batching inbound
	// messages instead of sending them out immediately to subscribers as per
	// the configured policy.
	BatchingEnabled bool `json:"batchingEnabled"`
	// When configured, sets the maximium batching interval in the channel.
	BatchingInterval *int `json:"batchingInterval,omitempty"`
	// If `true`, enables conflation for channels within this namespace.
	// Conflation reduces the number of messages sent to subscribers by
	// combining multiple messages into a single message.
	ConflationEnabled bool `json:"conflationEnabled"`
	// The interval in milliseconds at which messages are conflated. This
	// determines how frequently messages are combined into a single message.
	ConflationInterval *int `json:"conflationInterval"`
	// The key used to determine which messages should be conflated. Messages
	// with the same conflation key will be combined into a single message.
	ConflationKey string `json:"conflationKey"`
}

// Namespaces lists the namespaces for the specified application ID.
func (c *Client) Namespaces(appID string) ([]Namespace, error) {
	var namespaces []Namespace
	err := c.request("GET", "/apps/"+appID+"/namespaces", nil, &namespaces)
	return namespaces, err
}

// CreateNamespace creates a namespace for the specified application ID.
func (c *Client) CreateNamespace(appID string, namespace *Namespace) (Namespace, error) {
	var out Namespace
	err := c.request("POST", "/apps/"+appID+"/namespaces", &namespace, &out)
	return out, err
}

// UpdateNamespace updates the namespace with the specified ID, for the application with the specified application ID.
func (c *Client) UpdateNamespace(appID string, namespace *Namespace) (Namespace, error) {
	in := *namespace
	id := in.ID
	in.ID = ""

	var out Namespace
	err := c.request("PATCH", "/apps/"+appID+"/namespaces/"+id, &in, &out)
	return out, err
}

// DeleteNamespace deletes the namespace with the specified ID, for the specified application ID.
func (c *Client) DeleteNamespace(appID, namespaceID string) error {
	err := c.request("DELETE", "/apps/"+appID+"/namespaces/"+namespaceID, nil, nil)
	return err
}

func Interval(interval int) *int {
	return &interval
}
