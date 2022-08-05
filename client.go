// package control is an implementation of the Ably Control API.
//
// The Ably Control API is a REST API that enables you to manage your Ably
// account programmatically https://ably.com/documentation/control-api.
//
// Using the Control API you can automate the provisioning, management,
// and testing of your Ably realtime infrastructure. You can dynamically
// create Ably apps, configure them, and delete them if necessary.
//
// With the Control API you can create and manage:
//   - Your Ably apps
//   - API keys for an Ably app
//   - Namespaces (for channel rules)
//   - Queues
//   - Integration rules
//
// Control API is currently in Preview.
package control

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// The URL of the Ably Control API.
const API_URL = "https://control.ably.net/v1"

// Client represents a REST client for the Ably Control API.
type Client struct {
	token     string
	accountID string
	// Url is the base url for the REST API.
	Url string
}

// NewClient creates a new REST client.
//
// Creating a new client involves making a request to the REST API to
// fetch the account ID accociated with the token.
func NewClient(token string) (Client, Me, error) {
	return NewClientWithURL(token, API_URL)
}

// / NewClientWithURL is the same as NewClient but also takes a custom url.
func NewClientWithURL(token, url string) (Client, Me, error) {
	client := Client{
		token: token,
		Url:   url,
	}
	me, err := client.Me()
	if err != nil {
		return client, me, err
	}
	client.accountID = me.Account.ID
	return client, me, nil
}

func (c *Client) request(method, path string, in, out interface{}) error {
	var inR io.Reader
	if in != nil {
		inData, err := json.Marshal(in)
		if err != nil {
			return err
		}
		inR = bytes.NewReader(inData)
	}
	req, err := http.NewRequest(method, c.Url+path, inR)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		var errorInfo ErrorInfo
		err = json.Unmarshal(body, &errorInfo)
		if err == nil {
			return errorInfo
		} else {
			return ErrorInfo{
				Message:    string(body),
				Code:       0,
				StatusCode: res.StatusCode,
				HRef:       "",
			}
		}
	}
	if out != nil {
		return json.NewDecoder(res.Body).Decode(out)
	}
	return nil
}
