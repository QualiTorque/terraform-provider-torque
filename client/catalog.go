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
	fmt.Println(&data)
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall label update request: %s", err)
	}
	fmt.Println(string(payload))
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
