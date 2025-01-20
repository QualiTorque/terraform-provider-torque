package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateDeploymentEngine(engine_type string, name string, description string, agent_name string, auth_token string, polling_interval_seconds int32, server_url string, allowed_spaces AllowedSpaces) error {
	data := DeploymentEngine{
		Name:                   name,
		Description:            description,
		Type:                   engine_type,
		ServerUrl:              server_url,
		AgentName:              agent_name,
		AuthToken:              auth_token,
		PollingIntervalSeconds: polling_interval_seconds,
		AllowedSpaces:          allowed_spaces,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall deployment engine request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/deployment_engines", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) GetDeploymentEngine(name string) (*DeploymentEngineRead, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/deployment_engines/%s", c.HostURL, name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	deployment_engine := DeploymentEngineRead{}

	err = json.Unmarshal(body, &deployment_engine)
	if err != nil {
		return nil, err
	}
	return &deployment_engine, nil
}

func (c *Client) UpdateDeploymentEngine(engine_type string, current_name string, name string, description string, agent_name string, auth_token string, polling_interval_seconds int32, server_url string, allowed_spaces AllowedSpaces) error {
	data := DeploymentEngine{
		Name:                   name,
		Description:            description,
		Type:                   engine_type,
		ServerUrl:              server_url,
		AgentName:              agent_name,
		AuthToken:              auth_token,
		PollingIntervalSeconds: polling_interval_seconds,
		AllowedSpaces:          allowed_spaces,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall deployment engine request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/deployment_engines/%s", c.HostURL, current_name), bytes.NewReader(payload))
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

func (c *Client) DeleteDeploymentEngine(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/deployment_engines/%s", c.HostURL, name), nil)
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
