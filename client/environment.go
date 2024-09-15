package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetEnvironmentDetails(spaceName string, environmentId string) (*Environment, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/environments/%s", c.HostURL, spaceName, environmentId), nil)
	if err != nil {
		return nil, "", err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, "", err
	}
	rawJSON := string(body)

	environment := Environment{}
	err = json.Unmarshal(body, &environment)
	if err != nil {
		return nil, rawJSON, err
	}

	return &environment, rawJSON, nil
}
