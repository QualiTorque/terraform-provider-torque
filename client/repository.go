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

func (c *Client) OnboardGitlabEnterpriseRepoToSpace(space_name string, repository_name string, repository_url string, token *string, branch string, credential_name string) error {
	data := GitlabEnterpriseRepoSpaceAssociation{
		Token:          token,
		Name:           repository_name,
		URL:            repository_url,
		Branch:         branch,
		CredentialName: credential_name,
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

func (c *Client) OnboardRepoToSpace(space_name string, repo_name string, repo_type string, repo_url string, repo_token string, repo_branch string) error {
	data := RepoSpaceAssociation{
		Type:        repo_type,
		URL:         repo_url,
		AccessToken: repo_token,
		Branch:      repo_branch,
		Name:        repo_name,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/repositories", c.HostURL, space_name), bytes.NewReader(payload))
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
