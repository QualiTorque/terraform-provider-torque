package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateAuditTarget(audit_type string, properties *AuditProperties) error {
	data := Audit{
		Type:       audit_type,
		Properties: properties,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall create audit target request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/settings/audit/config", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) GetAudit() (*Audit, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/settings/audit/config", c.HostURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	audit := Audit{}

	err = json.Unmarshal(body, &audit)
	if err != nil {
		return nil, err
	}
	return &audit, nil
}

func (c *Client) DeleteAudit(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/settings/audit/config/%s", c.HostURL, name), nil)
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
