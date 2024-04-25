package client

type User struct {
	Email                string   `json:"email"`
	FirstName            string   `json:"first_name"`
	LastName             string   `json:"last_name"`
	Timezone             string   `json:"timezone"`
	DisplayFirstName     string   `json:"display_first_name"`
	DisplayLastName      string   `json:"display_last_name"`
	UserType             string   `json:"user_type"`
	JoinDate             string   `json:"join_date"`
	AccountRole          string   `json:"account_role"`
	HasAccessToAllSpaces bool     `json:"has_access_to_all_spaces"`
	Permissions          []string `json:"permissions"`
}

type Blueprint struct {
	BlueprintName string `json:"blueprint_name"`
	Name          string `json:"name"`
	DisplayName   string `json:"display_name"`
	RepoName      string `json:"repository_name"`
	RepoBranch    string `json:"repository_branch"`
	Commit        string `json:"commit"`
	Description   string `json:"description"`
	Url           string `json:"url"`
	ModifiedBy    string `json:"modified_by"`
	Published     bool   `json:"enabled"`
}

type Environment struct {
	EnvironmentName string            `json:"environment_name"`
	BlueprintName   string            `json:"blueprint_name"`
	OwnerEmail      string            `json:"owner_email"`
	Description     string            `json:"description"`
	Inputs          map[string]string `json:"inputs"`
	// Tags            map[string]string `json:"tags"`
	// Collaborators   struct {
	// 	CollaboratorsEmails []string `json:"collaborators_emails"`
	// 	AllSpaceMembers     bool     `json:"all_space_members"`
	// } `json:"collaborators"`
	Automation bool `json:"automation"`
	// ScheduledEndTime string `json:"scheduled_end_time"`
	Duration string `json:"duration"`
	Id       string `json:"id"`
	// Source           struct {
	// 	BlueprintName  string `json:"blueprint_name"`
	// 	RepositoryName string `json:"repository_name"`
	// 	Branch         string `json:"branch"`
	// 	Commit         string `json:"commit"`
	// } `json:"source"`
	// Workflows []struct {
	// 	Name      string `json:"name"`
	// 	Schedules []struct {
	// 		Scheduler  string `json:"scheduler"`
	// 		Overridden bool   `json:"overridden"`
	// 	} `json:"schedules"`
	// 	Reminder        int               `json:"reminder"`
	// 	InputsOverrides map[string]string `json:"inputs_overrides"`
	// } `json:"workflows"`
}

type UserSpaceAssociation struct {
	Email     string `json:"email"`
	SpaceRole string `json:"space_role"`
}

type Space struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type AgentSpaceAssociation struct {
	Type                  string `json:"type"`
	DefaultNamespace      string `json:"namespace"`
	DefaultServiceAccount string `json:"service_account"`
}

type RepoSpaceAssociation struct {
	URL         string `json:"repository_url"`
	AccessToken string `json:"access_token"`
	Type        string `json:"repository_type"`
	Branch      string `json:"branch"`
	Name        string `json:"repository_name"`
}

type TagNameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Tag struct {
	Name           string   `json:"tag_name"`
	Value          string   `json:"tag_value"`
	Scope          string   `json:"tag_scope"`
	Description    string   `json:"description"`
	PossibleValues []string `json:"possible_values"`
}

type TagDetails struct {
	Name           string   `json:"name"`
	Value          string   `json:"value"`
	Scope          string   `json:"created_by"`
	Description    string   `json:"description"`
	PossibleValues []string `json:"possible_values"`
}

type CatalogItemRequest struct {
	BlueprintName  string `json:"blueprint_name"`
	RepositoryName string `json:"repository_name"`
}

type ParameterRequest struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Sensitive   bool   `json:"sensitive"`
	Description string `json:"description"`
}

type GroupRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	IdpId       string      `json:"idp_identifier"`
	Users       []string    `json:"users"`
	AccountRole string      `json:"account_role"`
	SpaceRoles  []SpaceRole `json:"space_roles"`
}

type SpaceRole struct {
	SpaceName string `json:"space_name"`
	SpaceRole string `json:"space_role"`
}

type AwsCostTaret struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	ARN        string `json:"role_arn"`
	ExternalId string `json:"external_id"`
}

type SubscriptionsRequest struct {
	Name                  string                                    `json:"name"`
	Description           string                                    `json:"description"`
	EndThreshold          int64                                     `json:"end_threshold"`
	WorkflowStartReminder int64                                     `json:"workflow_start_reminder"`
	Target                SubscriptionsTargetRequest                `json:"target"`
	Events                []string                                  `json:"event_types"`
	WorkflowEventNotifier SubscriptionsWorkflowEventNotifierRequest `json:"workflow_events_notifier"`
	IdleReminder          []ReminderRequest                         `json:"idle_reminders"`
}

type SubscriptionsTargetRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type SubscriptionsWorkflowEventNotifierRequest struct {
	NotifyOnAllWorkflows bool `json:"notify_on_all_workflows"`
}

type ReminderRequest struct {
	TimeInHours int64 `json:"time_in_hours"`
}

type Account struct {
	ParentAccount string `json:"parent_account"`
	AccountName   string `json:"account_name"`
	Password      string `json:"password"`
	Company       string `json:"company"`
}
