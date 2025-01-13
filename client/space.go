package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (c *Client) CreateSpace(name string, color string, icon string) error {
	space := Space{
		Name:  name,
		Color: color,
		Icon:  icon,
	}

	payload, err := json.Marshal(space)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) DeleteSpace(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s", c.HostURL, name), nil)
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

func (c *Client) AddAgentToSpace(agent string, ns string, sa string, space string, agnet_type string) error {
	fmt.Println(c.HostURL + "api/spaces")

	data := AgentSpaceAssociation{
		Type:                  agnet_type,
		DefaultNamespace:      ns,
		DefaultServiceAccount: sa,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall agent association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/agents/%s", c.HostURL, space, agent), bytes.NewReader(payload))
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

func (c *Client) RemoveAgentFromSpace(agent string, space string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/agents/%s", c.HostURL, space, agent), nil)
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

// /api/spaces/{space_name}/blueprints.
func (c *Client) GetSpaceBlueprints(space_name string) ([]Blueprint, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/blueprints", c.HostURL, space_name), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	blueprints := []Blueprint{}
	err = json.Unmarshal(body, &blueprints)
	if err != nil {
		return nil, err
	}

	fmt.Println("Blueprint list length: " + strconv.Itoa(len(blueprints)))

	return blueprints, nil
}

func (c *Client) CreateSpaceTagValue(space_name string, tag_name string, tag_value string) error {
	data := NameValuePair{
		Name:  tag_name,
		Value: tag_value,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall space tag key value association: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/settings/tags", c.HostURL, space_name), bytes.NewReader(payload))
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

func (c *Client) SetSpaceTagValue(space_name string, tag_name string, tag_value string) error {
	data := NameValuePair{
		Name:  tag_name,
		Value: tag_value,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall space tag key value association: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/settings/tags/%s", c.HostURL, space_name, tag_name), bytes.NewReader(payload))
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

func (c *Client) DeleteSpaceTagValue(space_name string, tag_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/settings/tags/%s", c.HostURL, space_name, tag_name), nil)
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

func (c *Client) DeleteBlueprintTagValue(space_name string, tag_name string, repository_name string, blueprint_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/repositories/%s/blueprints/%s/settings/tags/%s", c.HostURL, space_name, repository_name, blueprint_name, tag_name), nil)
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

func (c *Client) CreateBlueprintTagValue(space_name string, tag_name string, tag_value string, repo_name string, blueprint_name string) error {
	data := NameValuePair{
		Name:  tag_name,
		Value: tag_value,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall blueprint tag key value association: %s", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/repositories/%s/blueprints/%s/settings/tags", c.HostURL, space_name, repo_name, blueprint_name), bytes.NewReader(payload))
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

func (c *Client) SetBlueprintTagValue(space_name string, tag_name string, tag_value string, repo_name string, blueprint_name string) error {
	data := NameValuePair{
		Name:  tag_name,
		Value: tag_value,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall blueprint tag key value association: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/repositories/%s/blueprints/%s/settings/tags/%s", c.HostURL, space_name, repo_name, blueprint_name, tag_name), bytes.NewReader(payload))
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

func (c *Client) AddGroupToSpace(groupName string, description string, idpId string, users []string, accountRole string, spaceRole []SpaceRole) error {

	data := GroupRequest{
		Name:        groupName,
		Description: description,
		IdpId:       idpId,
		Users:       users,
		AccountRole: accountRole,
		SpaceRoles:  spaceRole,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall group request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/groups", c.HostURL), bytes.NewReader(payload))
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

func (c *Client) DeleteGroup(group_name string) error {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/groups/%s", c.HostURL, group_name), nil)
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

func (c *Client) GetGroup(group_name string) (GroupRequest, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/groups", c.HostURL), nil)
	if err != nil {
		return GroupRequest{}, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return GroupRequest{}, err
	}

	groups := []GroupRequest{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return GroupRequest{}, err
	}

	group := GroupRequest{}
	for _, n := range groups {
		if group_name == n.Name {
			group = n
		}
	}

	return group, nil
}

func (c *Client) GetSpace(space_name string) (Space, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s", c.HostURL, space_name), nil)
	if err != nil {
		return Space{}, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return Space{}, err
	}

	space := Space{}
	err = json.Unmarshal(body, &space)
	if err != nil {
		return Space{}, err
	}

	return space, nil
}

func (c *Client) GetSpaces() ([]Space, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	spaces := []Space{}
	err = json.Unmarshal(body, &spaces)
	if err != nil {
		return nil, err
	}

	return spaces, nil
}

func (c *Client) UpdateAccountTag(name string, value string, description string, possible_values []string, scope string) error {

	tag := Tag{
		Name:           name,
		Value:          value,
		Scope:          scope,
		Description:    description,
		PossibleValues: possible_values,
	}

	payload, err := json.Marshal(tag)
	if err != nil {
		log.Fatalf("impossible to marshall update group request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/settings/tags/%s", c.HostURL, name), bytes.NewReader(payload))
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

func (c *Client) UpdateGroup(groupName string, description string, idpId string, users []string, accountRole string, spaceRole []SpaceRole) error {

	data := GroupRequest{
		Name:        groupName,
		Description: description,
		IdpId:       idpId,
		Users:       users,
		AccountRole: accountRole,
		SpaceRoles:  spaceRole,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall update group request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/groups/%s", c.HostURL, groupName), bytes.NewReader(payload))
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

func (c *Client) UpdateSpace(current_space string, name string, color string, icon string) error {

	data := Space{
		Name:  name,
		Color: color,
		Icon:  icon,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall update space request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s", c.HostURL, current_space), bytes.NewReader(payload))
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

func (c *Client) CreateSpaceEmailNotification(space_name string, notification_name string, environment_launched bool,
	environment_deployed bool, environment_force_ended bool, environment_idle bool, environment_extended bool,
	drift_detected bool, workflow_failed bool, workflow_started bool, updates_detected bool,
	collaborator_added bool, action_failed bool, environment_ending_failed bool, environment_ended bool,
	environment_active_with_error bool, workflow_start_reminder int64, end_threshold int64, idle_reminder []int64) (string, error) {

	data := SubscriptionsRequest{
		Name:        notification_name,
		Description: "",
		Target: SubscriptionsTargetRequest{
			Type:        "Email",
			Description: "",
		},
	}

	if environment_launched {
		data.Events = append(data.Events, "EnvironmentLaunched")
	}
	if environment_deployed {
		data.Events = append(data.Events, "EnvironmentDeployed")
	}
	if environment_force_ended {
		data.Events = append(data.Events, "EnvironmentForceEnded")
	}

	if environment_idle {
		data.Events = append(data.Events, "EnvironmentIdle")
		data.IdleReminder = []ReminderRequest{}
		for _, idleNumber := range idle_reminder {
			data.IdleReminder = append(data.IdleReminder, ReminderRequest{TimeInHours: idleNumber})
		}
	}

	if environment_extended {
		data.Events = append(data.Events, "EnvironmentExtended")
	}
	if drift_detected {
		data.Events = append(data.Events, "DriftDetected")
	}
	if workflow_failed {
		data.Events = append(data.Events, "WorkflowFailed")
	}
	if workflow_started {
		data.Events = append(data.Events, "WorkflowStarted")
		data.WorkflowStartReminder = workflow_start_reminder
	}
	if updates_detected {
		data.Events = append(data.Events, "UpdatesDetected")
	}
	if collaborator_added {
		data.Events = append(data.Events, "CollaboratorAdded")
	}
	if action_failed {
		data.Events = append(data.Events, "ActionFailed")
	}
	if environment_ending_failed {
		data.Events = append(data.Events, "EnvironmentEndingFailed")
	}
	if environment_ended {
		data.Events = append(data.Events, "EnvironmentEnded")
		data.EndThreshold = end_threshold
	}
	if environment_active_with_error {
		data.Events = append(data.Events, "EnvironmentActiveWithError")
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall update space request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/subscriptions", c.HostURL, space_name), bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) DeleteSpaceNotification(space_name string, notification_id string) error {
	fmt.Println(c.HostURL + "api/spaces")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/subscriptions?subscriptionId=%s", c.HostURL, space_name, notification_id), nil)
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
