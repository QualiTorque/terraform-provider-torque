package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
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
	const (
		maxRetries = 5
		delay      = 2 * time.Second
		timeout    = 30 * time.Second
	)

	endTime := time.Now().Add(timeout)

	for retries := 0; retries < maxRetries; retries++ {
		if time.Now().After(endTime) {
			return nil, fmt.Errorf("timed out waiting for blueprint %s to be available in the asset library", blueprint_name)
		}

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

		for _, blueprint_item := range blueprints {
			if blueprint_name == blueprint_item.Name {
				return &blueprint_item, nil
			}
		}

		time.Sleep(delay)
	}

	return nil, fmt.Errorf("blueprint %s not found after %d retries", blueprint_name, maxRetries)
}
