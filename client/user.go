package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) GetUserDetails(userEmail string) (*User, error) {
	fmt.Println(c.HostURL + "api/accounts/users/" + userEmail)
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/accounts/users/%s", c.HostURL, userEmail), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) AddUserToSpace(userEmail string, role string, space string) error {
	fmt.Println(c.HostURL + "api/spaces")

	user := UserSpaceAssociation{
		Email:     userEmail,
		SpaceRole: role,
	}

	payload, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/users", c.HostURL, space), bytes.NewReader(payload))
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
