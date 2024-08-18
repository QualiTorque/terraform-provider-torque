package client

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) AddBlueprintToAssetLibrary(space_name string, repo_name *string, blueprint_name string) error {
	baseURL := fmt.Sprintf("%sapi/spaces/%s/asset-library/%s", c.HostURL, space_name, blueprint_name)

	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}

	if repo_name != nil && *repo_name != "" {
		q := u.Query()
		q.Add("repository_name", *repo_name)
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("POST", u.String(), nil)
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
