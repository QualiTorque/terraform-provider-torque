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
	credentials := GitCredentialsRequest{
		SpaceName:       space_name,
		Name:            name,
		Description:     description,
		CloudType:       cloudtype,
		CloudIdentifier: repo_type,
		CredentialData:  credential_data,
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
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
