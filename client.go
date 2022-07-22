package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const API_URL = "https://control.ably.net/v1"

type Client struct {
	token     string
	accountID string
	Url       string
}

func NewClient(token string) (Client, Me, error) {
	return NewClientWithURL(token, API_URL)
}

func NewClientWithURL(token, url string) (Client, Me, error) {
	client := Client{
		token: token,
		Url:   API_URL,
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
		return fmt.Errorf("unexpected HTTP status: %s: %s", res.Status, body)
	}
	if out != nil {
		return json.NewDecoder(res.Body).Decode(out)
	}
	return nil
}
