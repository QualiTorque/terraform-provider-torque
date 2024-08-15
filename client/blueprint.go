package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetBlueprintDetails(space_name string, name string) (DetailedBlueprint, error) {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/catalog/%s", c.HostURL, space_name, name), nil)
	if err != nil {
		return DetailedBlueprint{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return DetailedBlueprint{}, err
	}

	blueprint := DetailedBlueprint{}
	err = json.Unmarshal(body, &blueprint)
	if err != nil {
		return DetailedBlueprint{}, err
	}
	return blueprint, nil
}
