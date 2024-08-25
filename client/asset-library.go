package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) AddBlueprintToAssetLibrary(space_name string, repo_name *string, blueprint_name string) error {
	baseURL := fmt.Sprintf("%sapi/spaces/%s/asset-library/%s", c.HostURL, space_name, blueprint_name)

	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	if repo_name != nil && *repo_name != "" {
		q := u.Query()
		q.Add("repository_name", *repo_name)
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("POST", u.String(), nil)
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

func (c *Client) RemoveBlueprintFromAssetLibrary(space_name string, repo_name *string, blueprint_name string) error {
	baseURL := fmt.Sprintf("%sapi/spaces/%s/asset-library/%s", c.HostURL, space_name, blueprint_name)

	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	if repo_name != nil && *repo_name != "" {
		q := u.Query()
		q.Add("repository_name", *repo_name)
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
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

func (c *Client) GetBlueprintFromAssetLibrary(space_name string, blueprint_name string) (*Blueprint, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/asset-library", c.HostURL, space_name), nil)
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
		if blueprint_name == blueprint_item.Name {
			blueprint = blueprint_item
			return &blueprint, nil
		}
	}

	return nil, fmt.Errorf("blueprint %s not found", blueprint_name)
}
