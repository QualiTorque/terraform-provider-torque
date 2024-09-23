package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) AddTag(name string, value string, description string, possible_values []string, scope string) error {
	tag := Tag{
		Name:           name,
		Value:          value,
		Scope:          scope,
		Description:    description,
		PossibleValues: possible_values,
	}
	payload, err := json.Marshal(tag)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/settings/tags", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) GetTag(tag_name string) (*Tag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/settings/tags", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	tags := []Tag{}
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}
	tag := Tag{}
	for _, tag_item := range tags {
		if tag_name == tag_item.Name {
			tag = tag_item
			return &tag, nil
		}
	}
	return nil, fmt.Errorf("tag %s not found", tag_name)
}

func (c *Client) UpdateTag(current_name string, name string, value string, description string, possible_values []string, scope string) error {
	tag := Tag{
		Name:           name,
		Value:          value,
		Scope:          scope,
		Description:    description,
		PossibleValues: possible_values,
	}

	payload, err := json.Marshal(tag)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/settings/tags/%s", c.HostURL, current_name), bytes.NewReader(payload))
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

func (c *Client) RemoveTag(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/settings/tags/%s", c.HostURL, name), nil)
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

func (c *Client) GetSpaceTags(space_name string) ([]Tag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/settings/tags", c.HostURL, space_name), nil)
	if err != nil {
		return []Tag{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return []Tag{}, err
	}

	blueprint := []Tag{}
	err = json.Unmarshal(body, &blueprint)
	if err != nil {
		return []Tag{}, err
	}
	return blueprint, nil
}

func (c *Client) GetSpaceTag(space_name string, tag_name string) (NameValuePair, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/settings/tags", c.HostURL, space_name), nil)
	if err != nil {
		return NameValuePair{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return NameValuePair{}, err
	}
	tags := []NameValuePair{}
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return NameValuePair{}, err
	}
	// var tag NameValuePair
	for _, tag := range tags {
		if tag.Name == tag_name {
			return tag, nil
		}
	}
	return NameValuePair{}, fmt.Errorf("Tag '%s' not found in space '%s'", tag_name, space_name)
}
