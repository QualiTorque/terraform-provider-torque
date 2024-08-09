package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetIntrospectionDetails(spaceName string, environmentId string) ([]IntrospectionItem, error) {
	url := fmt.Sprintf("%sapi/spaces/%s/environments/%s/introspection", c.HostURL, spaceName, environmentId)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var introspection []IntrospectionItem
	err = json.Unmarshal(body, &introspection)
	if err != nil {
		return nil, err
	}

	// Return the introspection details
	return introspection, nil
}
