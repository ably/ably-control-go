package control

// A struct representing an Ably Key.
type Key struct {
	// The key ID.
	ID string `json:"id,omitempty"`
	// The Ably application ID which this key is associated with.
	AppID string `json:"appId,omitempty"`
	// The name for your API key. This is a friendly name for your reference.
	Name string `json:"name,omitempty"`
	// The status of the key. 0 is enabled, 1 is revoked.
	Status int `json:"status"`
	// The complete API key including API secret.
	Key string `json:"key,omitempty"`
	// The capabilities that this key has. More information on capabilities
	// can be found in the Ably documentation https://ably.com/documentation/core-features/authentication#capabilities-explained.
	Capability map[string][]string `json:"capability"`
	// Unix timestamp representing the date and time of creation of the key.
	Created int `json:"created"`
	// Unix timestamp representing the date and time of the last modification of the key.
	Modified int `json:"modified"`
	// Token revocation is a security mechanism allowing an app to invalidate authentication tokens,
	// primarily used against malicious clients. Implementation sets tokens' maximum time-to-live (TTL) to one hour.
	RevocableTokens bool `json:"revocableTokens"`
}

// A struct representing the settable fields of an Ably key.
type NewKey struct {
	// The name for your API key. This is a friendly name for your reference.
	Name string `json:"name,omitempty"`
	// The capabilities that this key has. More information on capabilities
	// can be found in the Ably documentation https://ably.com/documentation/core-features/authentication#capabilities-explained.
	Capability map[string][]string `json:"capability"`
	// Enable Revocable Tokens. More information on Token Revocation can be
	// found in the Ably documentation https://ably.com/docs/auth/revocation
	RevocableTokens bool `json:"revocableTokens"`
}

// Keys lists the API keys associated with the application ID.
func (c *Client) Keys(appID string) ([]Key, error) {
	var keys []Key
	err := c.request("GET", "/apps/"+appID+"/keys", nil, &keys)
	return keys, err
}

// CreateKey creates an application with the specified properties.
func (c *Client) CreateKey(appID string, key *NewKey) (Key, error) {
	var out Key
	err := c.request("POST", "/apps/"+appID+"/keys", &key, &out)
	return out, err
}

// UpdateKey updates the API key with the specified key ID.
func (c *Client) UpdateKey(appID, keyID string, key *NewKey) (Key, error) {
	var out Key
	err := c.request("PATCH", "/apps/"+appID+"/keys/"+keyID, &key, &out)
	return out, err
}

// RevokeKey revokes the API key with the specified ID. This deletes the key.
func (c *Client) RevokeKey(appID, keyID string) error {
	err := c.request("POST", "/apps/"+appID+"/keys/"+keyID+"/revoke", nil, nil)
	return err
}
