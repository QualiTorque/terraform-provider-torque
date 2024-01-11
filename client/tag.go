package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) AddTag(name string, value string, description string, possible_values []string, scope string) error {
	fmt.Println(c.HostURL + "api/spaces")

	user := Tag{
		Name:           name,
		Value:          value,
		Scope:          scope,
		Description:    description,
		PossibleValues: possible_values,
	}

	payload, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/settings/tags", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) RemoveTag(name string) error {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/settings/tags/%s", c.HostURL, name), nil)
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
