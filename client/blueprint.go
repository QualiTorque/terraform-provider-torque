package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetBlueprint(space_name string, name string) (Blueprint, error) {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/catalog/%s", c.HostURL, space_name, name), nil)
	if err != nil {
		return Blueprint{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return Blueprint{}, err
	}

	blueprint := Blueprint{}
	err = json.Unmarshal(body, &blueprint)
	if err != nil {
		return Blueprint{}, err
	}
	return blueprint, nil
}
