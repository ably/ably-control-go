package control

// A struct representing an Ably application.
type App struct {
	//The application ID.
	ID string `json:"id,omitempty"`
	// The ID of your Ably account.
	AccountID string `json:"accountId,omitempty"`
	// The application name.
	Name string `json:"name,omitempty"`
	// The application status. Disabled applications will not accept
	// new connections and will return an error to all clients.
	Status string `json:"status,omitempty"`
	// Enforce TLS for all connections. This setting overrides any channel setting.
	TLSOnly bool `json:"tlsOnly"`
	// The Firebase Cloud Messaging key.
	FcmKey string `json:"fcmKey,omitempty"`
	// The Apple Push Notification service certificate.
	// This field can only be used to set a new value,
	// it will not be populated by queries.
	ApnsCertificate string `json:"apnsCertificate,omitempty"`
	// The Apple Push Notification service private key.
	// This field can only be used to set a new value,
	// it will not be populated by queries.
	ApnsPrivateKey string `json:"apnsPrivateKey,omitempty"`
	// The Apple Push Notification service sandbox endpoint.
	ApnsUseSandboxEndpoint bool `json:"apnsUseSandboxEndpoint"`
}

// Apps fetches a list of all your Ably apps.
func (c *Client) Apps() ([]App, error) {
	var apps []App
	err := c.request("GET", "/accounts/"+c.accountID+"/apps", nil, &apps)
	return apps, err
}

// CreateApp creates a new Ably app.
func (c *Client) CreateApp(app *App) (App, error) {
	var out App
	err := c.request("POST", "/accounts/"+c.accountID+"/apps", app, &out)
	return out, err
}

// UpdateApp updates an existing Ably app.
func (c *Client) UpdateApp(id string, app *App) (App, error) {
	var out App
	err := c.request("PATCH", "/apps/"+id, app, &out)
	return out, err
}

// DeleteApp deletes an Ably app.
func (c *Client) DeleteApp(id string) error {
	err := c.request("DELETE", "/apps/"+id, nil, nil)
	return err
}
