package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) GetEnvironmentDetails(spaceName string, environmentId string) (*Environment, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/environments/%s", c.HostURL, spaceName, environmentId), nil)
	if err != nil {
		return nil, "", err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, "", err
	}
	rawJSON := string(body)

	environment := Environment{}
	err = json.Unmarshal(body, &environment)
	if err != nil {
		return nil, rawJSON, err
	}

	return &environment, rawJSON, nil
}

func (c *Client) CreateEnvironment(Space string, BlueprintName string, EnvironmentName string, Duration string, Description string,
	Inputs map[string]string, OwnerEmail string, Automation bool, Tags map[string]string, Collaborators Collaborators, ScheduledEndTime string, BlueprintSource BlueprintSource, Workflows []EnvironmentWorkflow) ([]byte, error) {
	fmt.Println(c.HostURL + "api/spaces/" + Space + "/environments")

	environment := EnvironmentRequest{
		BlueprintName:    BlueprintName,
		EnvironmentName:  EnvironmentName,
		Description:      Description,
		Duration:         Duration,
		Inputs:           Inputs,
		OwnerEmail:       OwnerEmail,
		Automation:       Automation,
		Tags:             Tags,
		Collaborators:    Collaborators,
		ScheduledEndTime: ScheduledEndTime,
		BlueprintSource:  BlueprintSource,
		Workflows:        Workflows,
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

func (c *Client) UpdateEnvironmentName(Space string, Id string, Name string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/environments/%s/update_v2/%s/rename", c.HostURL, Space, Id, Name), nil)
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

func (c *Client) UpdateEnvironmentCollaborators(Space string, Id string, CollaboratorsEmails []string, AllSpaceMembers bool) error {
	collaborators := Collaborators{
		Collaborators:   CollaboratorsEmails,
		AllSpaceMembers: AllSpaceMembers,
	}
	payload, err := json.Marshal(collaborators)
	if err != nil {
		log.Fatalf("impossible to marshall Environment: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/environments/%s/collaborators", c.HostURL, Space, Id), bytes.NewReader(payload))
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

func (c *Client) TerminateEnvironment(Space string, Id string) error {
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

func (c *Client) ForceTerminateEnvironment(Space string, Id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/environments/force/%s", c.HostURL, Space, Id), nil)
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
