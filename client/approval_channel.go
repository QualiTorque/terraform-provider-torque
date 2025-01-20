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
