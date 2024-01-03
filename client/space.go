package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateSpace(name string, color string, icon string) error {
	fmt.Println(c.HostURL + "api/spaces")

	space := Space{
		Name:  name,
		Color: color,
		Icon:  icon,
	}

	payload, err := json.Marshal(space)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) DeleteSpace(name string) error {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s", c.HostURL, name), nil)
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

func (c *Client) AddAgentToSpace(agent string, ns string, sa string, space string, agnet_type string) error {
	fmt.Println(c.HostURL + "api/spaces")

	data := AgentSpaceAssociation{
		Type:                  agnet_type,
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

func (c *Client) OnboardRepoToSpace(space_name string, repo_name string, repo_type string, repo_url string, repo_token string, repo_branch string) error {
	fmt.Println(c.HostURL + "api/spaces")

	data := RepoSpaceAssociation{
		Type:        repo_type,
		URL:         repo_url,
		AccessToken: repo_token,
		Branch:      repo_branch,
		Name:        repo_name,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/repositories", c.HostURL, space_name), bytes.NewReader(payload))
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
