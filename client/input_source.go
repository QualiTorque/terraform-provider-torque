package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateInputSource(Name string, Description string, AllowedSpaces AllowedSpaces, Details InputSourceDetails) error {
	data := TorqueInputSource{
		Name:          Name,
		Description:   Description,
		AllowedSpaces: AllowedSpaces,
		Details:       Details,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall Input Source request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/input_sources", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) DeleteInputSource(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/input_sources/%s", c.HostURL, name), nil)
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

func (c *Client) UpdateInputSource(CurrentName string, Name string, Description string, AllowedSpaces AllowedSpaces, Details InputSourceDetails) error {
	data := TorqueInputSource{
		Name:          Name,
		Description:   Description,
		AllowedSpaces: AllowedSpaces,
		Details:       Details,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall Input Source request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/input_sources/%s", c.HostURL, CurrentName), bytes.NewReader(payload))
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

func (c *Client) GetInputSource(Name string) (*TorqueInputSource, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/input_sources/%s", c.HostURL, Name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	input_source := TorqueInputSource{}

	err = json.Unmarshal(body, &input_source)
	if err != nil {
		return nil, err
	}
	return &input_source, nil
}
