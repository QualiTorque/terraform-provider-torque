package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) PublishBlueprintInSpace(space_name string, repo_name string, blueprint_name string) error {
	data := CatalogItemRequest{
		BlueprintName:  blueprint_name,
		RepositoryName: repo_name,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/catalog", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) UnpublishBlueprintInSpace(space_name string, repo_name string, blueprint_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/catalog/%s?repository_name=%s", c.HostURL, space_name, blueprint_name, repo_name), nil)
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

func (c *Client) EditCatalogItemLabels(space_name string, blueprint_name string, repository_name string, labels []string) error {
	data := CatalogItemLabelsRequest{
		BlueprintName:  blueprint_name,
		RepositoryName: repository_name,
		Labels:         labels,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall label update request: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/catalog/labels", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) AllowLaunch(blueprint_name string, repository_name string, space_name string, launch_allowed bool) error {
	data := WorkflowRequest{
		BlueprintName:  blueprint_name,
		RepositoryName: repository_name,
		SpaceName:      space_name,
		LaunchAllowed:  launch_allowed,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall workflow request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/catalog/launch_allowed", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) SetCatalogItemCustomIcon(space_name string, blueprint_name string, repository_name string, key string) error {
	type setCatalogItemCustomIconRequest struct {
		BlueprintName  string `json:"blueprint_name"`
		RepositoryName string `json:"repository_name"`
		CustomIconKey  string `json:"custom_icon_key"`
	}
	request := setCatalogItemCustomIconRequest{
		BlueprintName:  blueprint_name,
		RepositoryName: repository_name,
		CustomIconKey:  key,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("impossible to marshall custom icon request: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/catalog/icons", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) SetCatalogItemIcon(space_name string, blueprint_name string, repository_name string, icon string) error {
	type setCatalogItemCustomIconRequest struct {
		BlueprintName  string `json:"blueprint_name"`
		RepositoryName string `json:"repository_name"`
		Icon           string `json:"icon"`
	}
	request := setCatalogItemCustomIconRequest{
		BlueprintName:  blueprint_name,
		RepositoryName: repository_name,
		Icon:           icon,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("impossible to marshall custom icon request: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/catalog/icons", c.HostURL, space_name), bytes.NewReader(payload))
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
