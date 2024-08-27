package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateEnvironmentLabel(key string, value string) error {
	data := KeyValuePair{
		Key:   key,
		Value: value,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall create environment label request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/environments/labels", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) GetEnvironmentLabel(key string) (*KeyValuePair, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/environments/labels", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	labels := []KeyValuePair{}
	err = json.Unmarshal(body, &labels)
	if err != nil {
		return nil, err
	}

	label := KeyValuePair{}
	for _, label_item := range labels {
		if key == label_item.Key {
			label = label_item
			return &label, nil
		}
	}

	return &label, nil
}

func (c *Client) UpdateEnvironmentLabel(current_key string, current_value, key string, value string) error {
	data := KeyValuePair{
		Key:   key,
		Value: value,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall label update request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/environments/labels/%s?label_value=%s", c.HostURL, current_key, current_value), bytes.NewReader(payload))
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

func (c *Client) DeleteEnvironmentLabel(key string, value string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/environments/labels/%s?label_value=%s", c.HostURL, key, value), nil)
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

func (c *Client) UpdateEnvironmentLabels(environment_id string, space_name string, added_labels []KeyValuePair, removed_labels []KeyValuePair) error {
	data := EnvironmentLabelsUpdateRequest{
		SpaceName:     space_name,
		EnvironmentId: environment_id,
		AddedLabels:   added_labels,
		RemovedLabels: removed_labels,
	}
	fmt.Println(data)
	payload, err := json.Marshal(data)
	fmt.Println(string(payload))
	if err != nil {
		log.Fatalf("impossible to marshall label update request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/environments/%s/labels", c.HostURL, space_name, environment_id), bytes.NewReader(payload))
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
