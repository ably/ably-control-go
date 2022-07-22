package control

type Namespace struct {
	ID               string `json:"id,omitempty"`
	Authenticated    bool   `json:"authenticated"`
	Persisted        bool   `json:"persisted"`
	PersistLast      bool   `json:"persistLast"`
	PushEnabled      bool   `json:"pushEnabled"`
	TlsOnly          bool   `json:"tlsOnly"`
	ExposeTimeserial bool   `json:"exposeTimeserial"`
}

func (c *Client) Namespaces(appID string) ([]Namespace, error) {
	var namespaces []Namespace
	err := c.request("GET", "/apps/"+appID+"/namespaces", nil, &namespaces)
	return namespaces, err
}

func (c *Client) CreateNamespace(appID string, namespace Namespace) (Namespace, error) {
	var out Namespace
	err := c.request("POST", "/apps/"+appID+"/namespaces", &namespace, &out)
	return out, err
}

func (c *Client) UpdateNamespace(appID, namespaceID string, namespace Namespace) (Namespace, error) {
	var out Namespace
	err := c.request("PATCH", "/apps/"+appID+"/namespaces/"+namespaceID, &namespace, &out)
	return out, err
}

func (c *Client) DeleteNamespace(appID, namespaceID string) error {
	err := c.request("DELETE", "/apps/"+appID+"/namespaces/"+namespaceID, nil, nil)
	return err
}
