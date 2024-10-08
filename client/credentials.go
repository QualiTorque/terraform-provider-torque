package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateSpaceGitCredentials(space_name string, name string, description string, repo_type string, token string) error {
	const cloudtype = "sourceControl"
	credential_data := CredentialData{
		Token: token,
		Type:  repo_type,
	}
	credentials := GitCredentials{
		SpaceName:       space_name,
		Name:            name,
		Description:     description,
		CloudType:       cloudtype,
		CloudIdentifier: repo_type,
		CredentialData:  credential_data,
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("impossible to marshall credentials: %s", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/settings/credentialstore", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) DeleteSpaceGitCredentials(space_name string, credential_name string) error {
	type DeleteSpaceCredentialRequest struct {
		SpaceName      string `json:"space_name"`
		CredentialName string `json:"credential_name"`
	}
	request := DeleteSpaceCredentialRequest{
		SpaceName:      space_name,
		CredentialName: credential_name,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("impossible to marshall credentials: %s", err)
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/settings/credentialstore/%s", c.HostURL, space_name, credential_name), bytes.NewReader(payload))
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

func (c *Client) GetSpaceGitCredentials(space_name string, credential_name string) (GitCredentials, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/settings/credentialstore/%s", c.HostURL, space_name, credential_name), nil)

	credentials := GitCredentials{}

	if err != nil {
		return credentials, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return credentials, err
	}

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		return credentials, err
	}

	return credentials, nil
}

func (c *Client) UpdateSpaceGitCredentials(space_name string, name string, description string, repo_type string, token string) error {
	const cloudtype = "sourceControl"
	credential_data := CredentialData{
		Token: token,
		Type:  repo_type,
	}
	credentials := GitCredentials{
		SpaceName:       space_name,
		Name:            name,
		Description:     description,
		CloudType:       cloudtype,
		CloudIdentifier: repo_type,
		CredentialData:  credential_data,
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("impossible to marshall credentials: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/settings/credentialstore/%s", c.HostURL, space_name, name), bytes.NewReader(payload))
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
