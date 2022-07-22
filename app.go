package control

type App struct {
	ID                     string `json:"id,omitempty"`
	AccountID              string `json:"accountId,omitempty"`
	Name                   string `json:"name,omitempty"`
	Status                 string `json:"status,omitempty"`
	TLSOnly                bool   `json:"tlsOnly"`
	FcmKey                 string `json:"fcmKey,omitempty"`
	ApnsCertificate        string `json:"apnsCertificate,omitempty"`
	ApnsPrivateKey         string `json:"apnsPrivateKey,omitempty"`
	ApnsUseSandboxEndpoint bool   `json:"apnsUseSandboxEndpoint"`
}

func (c *Client) Apps() ([]App, error) {
	var apps []App
	err := c.request("GET", "/accounts/"+c.accountID+"/apps", nil, &apps)
	return apps, err
}

func (c *Client) CreateApp(app *App) (App, error) {
	var out App
	err := c.request("POST", "/accounts/"+c.accountID+"/apps", app, &out)
	return out, err
}

func (c *Client) UpdateApp(id string, app *App) (App, error) {
	var out App
	err := c.request("PATCH", "/apps/"+id, app, &out)
	return out, err
}

func (c *Client) DeleteApp(id string) error {
	err := c.request("DELETE", "/apps/"+id, nil, nil)
	return err
}
