package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateSpaceCredentials(space_name string, name string, description string, cloudtype string, cloud_identifier string, token *string) error {
	credential_data := CredentialData{
		Token: token,
		Type:  cloud_identifier,
	}
	credentials := SpaceCredentials{
		SpaceName:       space_name,
		Name:            name,
		Description:     description,
		CloudType:       cloudtype,
		CloudIdentifier: cloud_identifier,
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

func (c *Client) CreateAccountCredentials(name string, description string, cloud_type string, cloud_identifier string, credential_type string, token *string, key *string, secret *string, allowed_space_names []string) error {
	credential_data := CredentialData{
		Token:  token,
		Type:   credential_type,
		Key:    key,
		Secret: secret,
	}

	credentials := AccountCredentials{
		Name:              name,
		Description:       description,
		CloudType:         cloud_type,
		CloudIdentifier:   cloud_identifier,
		CredentialData:    credential_data,
		AllowedSpaceNames: allowed_space_names,
		AllSpacesAllowed:  len(allowed_space_names) == 0 || allowed_space_names == nil,
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("impossible to marshall credentials: %s", err)
	}

	// fmt.Println(credential_data)
	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/settings/credentialstore", c.HostURL), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	fmt.Printf("Payload being sent:\n%s\n", string(payload))
	fmt.Println(string(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	_, err = c.doRequest(req, &c.Token)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteSpaceCredentials(space_name string, credential_name string) error {
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

func (c *Client) DeleteAccountCredentials(credential_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/settings/credentialstore/%s", c.HostURL, credential_name), nil)
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

func (c *Client) GetSpaceCredentials(space_name string, credential_name string) (SpaceCredentials, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/settings/credentialstore/%s", c.HostURL, space_name, credential_name), nil)

	credentials := SpaceCredentials{}

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

func (c *Client) GetCredentials(credential_name string) (AccountCredentials, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/settings/credentialstore/%s", c.HostURL, credential_name), nil)

	credentials := AccountCredentials{}

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

func (c *Client) UpdateSpaceCredentials(space_name string, name string, description string, cloudtype string, cloud_identifier string, token *string) error {
	credential_data := CredentialData{
		Token: token,
		Type:  cloud_identifier,
	}
	credentials := SpaceCredentials{
		SpaceName:       space_name,
		Name:            name,
		Description:     description,
		CloudType:       cloudtype,
		CloudIdentifier: cloud_identifier,
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

func (c *Client) UpdateAccountCredentials(name string, description string, cloud_identifier string, cloudtype string, credential_type string, token *string, key *string, secret *string, allowed_space_names []string) error {
	credential_data := CredentialData{
		Token:  token,
		Type:   credential_type,
		Key:    key,
		Secret: secret,
	}

	credentials := AccountCredentials{
		Name:              name,
		Description:       description,
		CloudType:         cloudtype,
		CloudIdentifier:   cloud_identifier,
		CredentialData:    credential_data,
		AllowedSpaceNames: allowed_space_names,
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("impossible to marshall credentials: %s", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/settings/credentialstore/%s", c.HostURL, name), bytes.NewReader(payload))
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
