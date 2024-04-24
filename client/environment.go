package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateEnvironment(Space string, BlueprintName string, EnvironmentName string, Duration string) ([]byte, error) {
	fmt.Println(c.HostURL + "api/spaces/" + Space + "/environments")

	environment := Environment{
		BlueprintName:   BlueprintName,
		EnvironmentName: EnvironmentName,
		Duration:        Duration,
	}

	payload, err := json.Marshal(environment)
	if err != nil {
		log.Fatalf("impossible to marshall Environment: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/environments", c.HostURL, Space), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) TerminateEnvironment(Space string, Id string) error {
	fmt.Println(c.HostURL + "api/spaces/" + Space + "/environments/" + Id)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/environments/%s", c.HostURL, Space, Id), nil)
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
