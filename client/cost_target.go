package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) AddAWSCostTarget(name string, target_type string, role_arn string, external_id string) error {

	data := AwsCostTarget{
		Name:       name,
		Type:       target_type,
		ARN:        role_arn,
		ExternalId: external_id,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall aws cost target request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/settings/costtargets", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) DeleteCostTarget(target_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/settings/costtargets/%s", c.HostURL, target_name), nil)
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

func (c *Client) UpdateAWSCostTarget(target_name string, new_target_name string, target_type string, role_arn string, external_id string) error {
	data := AwsCostTarget{
		Name:       target_name,
		Type:       target_type,
		ARN:        role_arn,
		ExternalId: external_id,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall target name update request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/settings/costtargets/%s", c.HostURL, target_name), bytes.NewReader(payload))
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
