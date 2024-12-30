package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

func (c *Client) SetBlueprintPolicies(space_name string, repository_name string, name string, max_duration string, default_duration string, default_extend string, max_active_environments *int32, always_on bool) error {
	data := Policies{
		MaxDuration:           max_duration,
		DefaultDuration:       default_duration,
		DefaultExtend:         default_extend,
		MaxActiveEnvironments: max_active_environments,
		AlwaysOn:              always_on,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/repositories/%s/blueprints/%s/policies", c.HostURL, space_name, repository_name, name), bytes.NewReader(payload))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	_, err = c.doRequest(req, &c.Token)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBlueprintDisplayName(space_name string, repository_name string, name string, display_name string) error {
	data := BlueprintDisplayNameRequest{
		BlueprintName:  name,
		RepositoryName: repository_name,
		DisplayName:    display_name,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to blueprint display name request: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/blueprints/display_name", c.HostURL, space_name), bytes.NewReader(payload))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	_, err = c.doRequest(req, &c.Token)
	if err != nil {
		return err
	}

	return nil
}
