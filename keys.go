package control

type Key struct {
	ID         string              `json:"id,omitempty"`
	AppID      string              `json:"appId,omitempty"`
	Name       string              `json:"name,omitempty"`
	Status     int                 `json:"status"`
	Key        string              `json:"key,omitempty"`
	Capability map[string][]string `json:"capability"`
	Created    int                 `json:"created"`
	Modified   int                 `json:"modified"`
}

type NewKey struct {
	Name       string              `json:"name,omitempty"`
	Capability map[string][]string `json:"capability"`
}

func (c *Client) Keys(appID string) ([]Key, error) {
	var keys []Key
	err := c.request("GET", "/apps/"+appID+"/keys", nil, &keys)
	return keys, err
}

func (c *Client) CreateKey(appID string, key *NewKey) (Key, error) {
	var out Key
	err := c.request("POST", "/apps/"+appID+"/keys", &key, &out)
	return out, err
}

func (c *Client) UpdateKey(appID, keyID string, key *NewKey) (Key, error) {
	var out Key
	err := c.request("PATCH", "/apps/"+appID+"/keys/"+keyID, &key, &out)
	return out, err
}

func (c *Client) RevokeKey(appID, keyID string) error {
	err := c.request("POST", "/apps/"+appID+"/keys/"+keyID+"/revoke", nil, nil)
	return err
}
