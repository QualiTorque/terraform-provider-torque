package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateLabel(space_name string, name string, color string) error {
	data := Label{
		Name:  name,
		Color: color,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall create label request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/labels", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) GetLabel(space_name string, name string) (*Label, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/labels", c.HostURL, space_name), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	labels := []Label{}
	err = json.Unmarshal(body, &labels)
	if err != nil {
		return nil, err
	}

	label := Label{}
	for _, label_item := range labels {
		if name == label_item.Name {
			label = label_item
			return &label, nil
		}
	}

	return &label, nil
}

func (c *Client) UpdateLabel(original_name string, space_name string, name string, color string) error {
	data := LabelRequest{
		OriginalName: original_name,
		Name:         name,
		Color:        color,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall label update request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/labels/update", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) DeleteLabel(space_name string, name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/labels?name=%s", c.HostURL, space_name, name), nil)
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
