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
		log.Fatalf("impossible to marshall create approval channel request: %s", err)
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
