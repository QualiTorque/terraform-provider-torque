package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateResourceInventory(credentials string, configuration ResourcesInventoryConfiguration) error {
	data := ResourceInventory{
		Credentials:   credentials,
		Configuration: configuration,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall Resource Inventory: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/cloudresource", c.HostURL), bytes.NewReader(payload))
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
