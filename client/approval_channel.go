package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateApprovalChannel(name string, description string, details ApprovalChannelDetails) error {
	data := ApprovalChannel{
		Name:        name,
		Description: description,
		Details:     details,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall create approval channel request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/approval/channels", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) GetApprovalChannel(name string) (*ApprovalChannel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/approval/channels/%s", c.HostURL, name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	approval_channel := ApprovalChannel{}

	err = json.Unmarshal(body, &approval_channel)
	if err != nil {
		return nil, err
	}
	return &approval_channel, nil
}

func (c *Client) UpdateApprovalChannel(name string, description string, details ApprovalChannelDetails) error {
	data := ApprovalChannel{
		Name:        name,
		Description: description,
		Details:     details,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall create approval channel request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/approval/channels/%s", c.HostURL, name), bytes.NewReader(payload))
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

func (c *Client) DeleteApprovalChannel(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/approval/channels/%s", c.HostURL, name), nil)
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
