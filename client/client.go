package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Hashicups URL.
const HostURL string = "https://portal.qtorque.io/"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Space      string
}

func NewClient(host, space, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 120 * time.Second},
		HostURL:    HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	c.Space = *space
	c.Token = *token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := c.Token

	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "terraform-provider-torque")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
