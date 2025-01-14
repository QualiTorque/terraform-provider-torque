package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) AddAgentToSpace(agent string, ns string, sa string, space string, agent_type string) error {
	data := AgentSpaceAssociation{
		Type:                  agent_type,
		DefaultNamespace:      ns,
		DefaultServiceAccount: sa,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/agents/%s", c.HostURL, space, agent), bytes.NewReader(payload))
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

func (c *Client) RemoveAgentFromSpace(agent string, space string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/agents/%s", c.HostURL, space, agent), nil)
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

func (c *Client) UpdateAgentSpaceAssociation(agent string, ns string, sa string, space string, agent_type string) error {
	type agentSpaceAssociationUpdateRequest struct {
		Type                  string `json:"type"`
		DefaultNamespace      string `json:"default_namespace"`
		DefaultServiceAccount string `json:"default_service_account"`
	}
	data := agentSpaceAssociationUpdateRequest{
		Type:                  agent_type,
		DefaultNamespace:      ns,
		DefaultServiceAccount: sa,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/executionhosts/k8s/%s/spaces/%s", c.HostURL, agent, space), bytes.NewReader(payload))
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
