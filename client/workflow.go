package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetWorkflow(workflow_name string) (Workflow, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/automation/workflows/%s", c.HostURL, workflow_name), nil)
	workflow := Workflow{}

	if err != nil {
		return workflow, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return workflow, err
	}

	err = json.Unmarshal(body, &workflow)
	if err != nil {
		return workflow, err
	}

	return workflow, nil
}
