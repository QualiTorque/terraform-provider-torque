package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetBlueprint(space_name string, name string) (*Blueprint, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/blueprints", c.HostURL, space_name), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}
	blueprints := []Blueprint{}
	err = json.Unmarshal(body, &blueprints)
	if err != nil {
		return nil, err
	}

	blueprint := Blueprint{}
	for _, blueprint_item := range blueprints {
		if name == blueprint_item.Name {
			blueprint = blueprint_item
			return &blueprint, nil
		}
	}

	return nil, fmt.Errorf("blueprint %s not found", name)
}
