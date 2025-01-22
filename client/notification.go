package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateSpaceEmailNotification(space_name string, notification_name string, environment_launched bool,
	environment_deployed bool, environment_force_ended bool, environment_idle bool, environment_extended bool,
	drift_detected bool, workflow_failed bool, workflow_started bool, updates_detected bool,
	collaborator_added bool, action_failed bool, environment_ending_failed bool, environment_ended bool,
	environment_active_with_error bool, workflow_start_reminder int64, end_threshold int64, blueprint_published bool, blueprint_unpublished bool, idle_reminder []int64) (string, error) {

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
	if blueprint_published {
		data.Events = append(data.Events, "BlueprintPublished")
	}
	if blueprint_unpublished {
		data.Events = append(data.Events, "BlueprintUnpublished")
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall space notification request: %s", err)
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

func (c *Client) UpdateSpaceEmailNotification(notification_id string, space_name string, notification_name string, environment_launched bool,
	environment_deployed bool, environment_force_ended bool, environment_idle bool, environment_extended bool,
	drift_detected bool, workflow_failed bool, workflow_started bool, updates_detected bool,
	collaborator_added bool, action_failed bool, environment_ending_failed bool, environment_ended bool,
	environment_active_with_error bool, workflow_start_reminder int64, end_threshold int64, blueprint_published bool, blueprint_unpublished bool, idle_reminder []int64) (string, error) {

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
	if blueprint_published {
		data.Events = append(data.Events, "BlueprintPublished")
	}
	if blueprint_unpublished {
		data.Events = append(data.Events, "BlueprintUnpublished")
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall update space request: %s", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%sapi/spaces/%s/subscriptions?subscriptionId=%s", c.HostURL, space_name, notification_id), bytes.NewReader(payload))
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
