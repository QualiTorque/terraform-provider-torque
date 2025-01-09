package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) ConfigureResourveInventory(credentials string, details ResourceInventoryDetails) error {
	data := ResourceInventory{
		Credentials: credentials,
		Details:     details,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall Resource Inventory: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/cloudresource/configuration", c.HostURL), bytes.NewReader(payload))
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

// func (c *Client) UpdateResourceInventory(credentials string, details ResourceInventoryDetails) error {
// 	data := ResourceInventory{
// 		Credentials: credentials,
// 		Details:     details,
// 	}

// 	payload, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("impossible to marshall Resource Inventory: %s", err)
// 	}

// 	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/cloudresource/configuration", c.HostURL), bytes.NewReader(payload))
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Accept", "application/json")

// 	_, err = c.doRequest(req, &c.Token)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (c *Client) DeleteResourceInventory(credentials string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/settings/credentialstore/%s", c.HostURL, credentials), nil)
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
