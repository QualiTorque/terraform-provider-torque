package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) AddSpaceParameter(space_name string, name string, value string, sensitive bool, description string) error {
	data := ParameterRequest{
		Name:        name,
		Value:       value,
		Sensitive:   sensitive,
		Description: description,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/settings/parameters", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) DeleteSpaceParameter(space_name string, parameter_name string) error {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/settings/parameters/%s", c.HostURL, space_name, parameter_name), nil)
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

func (c *Client) AddAccountParameter(name string, value string, sensitive bool, description string) error {
	data := ParameterRequest{
		Name:        name,
		Value:       value,
		Sensitive:   sensitive,
		Description: description,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/settings/parameters", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) GetSpaceParameter(space_name string, parameter_name string) (ParameterRequest, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/settings/parameters", c.HostURL, space_name), nil)
	if err != nil {
		return ParameterRequest{}, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return ParameterRequest{}, err
	}

	params := []ParameterRequest{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		return ParameterRequest{}, err
	}

	param := ParameterRequest{}
	for _, n := range params {
		if parameter_name == n.Name {
			param = n
		}
	}

	return param, nil
}

func (c *Client) GetAccountParameter(parameter_name string) (*ParameterRequest, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/settings/parameters/%s", c.HostURL, parameter_name), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	param := ParameterRequest{}
	err = json.Unmarshal(body, &param)
	if err != nil {
		return nil, err
	}

	return &param, nil
}

func (c *Client) DeleteAccountParameter(parameter_name string) error {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/settings/parameters/%s", c.HostURL, parameter_name), nil)
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

func (c *Client) UpdateAccountParameter(name string, value string, sensitive bool, description string) error {

	data := ParameterRequest{
		Name:        name,
		Value:       value,
		Sensitive:   sensitive,
		Description: description,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall update parameter request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/settings/parameters/%s", c.HostURL, name), bytes.NewReader(payload))
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

func (c *Client) UpdateSpaceParameter(space_name string, name string, value string, sensitive bool, description string) error {

	data := SpaceParameterRequest{
		Value:       value,
		Sensitive:   sensitive,
		Description: description,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall update space parameter request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/settings/parameters/%s", c.HostURL, space_name, name), bytes.NewReader(payload))
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
