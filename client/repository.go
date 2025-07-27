package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) OnboardCodeCommitRepoToSpace(space_name string, repository_name string, role_arn string, repository_url string, aws_region string,
	repository_branch string, external_id string, git_username string, git_password string, credential_name string) error {

	data := CodeCommitRepoSpaceAssociation{
		URL:            repository_url,
		RoleArn:        role_arn,
		Region:         aws_region,
		Branch:         repository_branch,
		Name:           repository_name,
		ExternalId:     external_id,
		Username:       git_username,
		Password:       git_password,
		CredentialName: credential_name,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/repositories/codeCommit", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) OnboardGitlabEnterpriseRepoToSpace(space_name string, repository_name string, repository_url string, token *string, branch string, credential_name string, agents []string, use_all_agents bool, auto_register_eac bool) error {
	data := GitlabEnterpriseRepoSpaceAssociation{
		Token:           token,
		Name:            repository_name,
		URL:             repository_url,
		Branch:          branch,
		CredentialName:  credential_name,
		Agents:          agents,
		UseAllAgents:    use_all_agents,
		AutoRegisterEac: auto_register_eac,
	}

	payload, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/repositories/gitlabEnterprise", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) OnboardAdoServerRepoToSpace(space_name string, repository_name string, repository_url string, token *string, branch string, credential_name string, agents []string, use_all_agents bool, auto_register_eac bool) error {
	data := AdoServerRepoSpaceAssociation{
		Token:           token,
		Name:            repository_name,
		URL:             repository_url,
		Branch:          branch,
		CredentialName:  credential_name,
		Agents:          agents,
		UseAllAgents:    use_all_agents,
		AutoRegisterEac: auto_register_eac,
	}

	payload, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/repositories/azureEnterprise", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) OnboardRepoToSpace(space_name string, repo_name string, repo_type string, repo_url string, repo_token *string, repo_branch string, credential_name *string) error {
	var data interface{}
	var url string
	if credential_name == nil || *credential_name == "" {
		data = RepoSpaceAssociation{
			URL:         repo_url,
			AccessToken: repo_token,
			Type:        repo_type,
			Branch:      repo_branch,
			Name:        repo_name,
		}
		url = fmt.Sprintf("%sapi/spaces/%s/repositories", c.HostURL, space_name)

	} else {
		data = RepoSpaceAssociationWithCredentials{
			URL:            repo_url,
			Type:           repo_type,
			Branch:         repo_branch,
			Name:           repo_name,
			CredentialName: credential_name,
		}
		url = fmt.Sprintf("%sapi/spaces/%s/repositories/%s", c.HostURL, space_name, repo_type)
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall repo association: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
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

func (c *Client) RemoveRepoFromSpace(space_name string, repo_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/repositories?repository_name=%s", c.HostURL, space_name, repo_name), nil)
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

func (c *Client) UpdateRepoCredentials(space_name string, repo_name string, credential_name string) error {
	data := map[string]string{
		"credential_name": credential_name,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall repo association: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/repositories/%s", c.HostURL, space_name, repo_name), bytes.NewReader(payload))
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

func (c *Client) UpdateRepoConfiguration(space_name string, repo_name string, credential_name string, agents []string, use_all_agents bool) error {
	data := RepoUpdate{
		CredentialName: credential_name,
		Agents:         agents,
		UseAllAgents:   use_all_agents,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall repo association: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/repositories/%s", c.HostURL, space_name, repo_name), bytes.NewReader(payload))
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

func (c *Client) GetRepoDetails(space_name string, repo_name string) (*RepoDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/repositories", c.HostURL, space_name), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}
	repos := []RepoDetails{}
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, err
	}
	repo := RepoDetails{}
	for _, repo_item := range repos {
		if repo_name == repo_item.Name {
			repo = repo_item
			return &repo, nil
		}
	}
	return nil, fmt.Errorf("repository %s not found", repo_name)
}
