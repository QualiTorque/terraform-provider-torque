package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetEnvironmentDetails(spaceName string, environmentId string) (*Environment, error) {
	fmt.Printf("%sapi/spaces/%s/environments/%s", c.HostURL, spaceName, environmentId)
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/environments/%s", c.HostURL, spaceName, environmentId), nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	environment := Environment{}
	err = json.Unmarshal(body, &environment)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}
