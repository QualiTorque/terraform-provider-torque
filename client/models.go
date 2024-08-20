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
	BlueprintName           string         `json:"blueprint_name"`
	Name                    string         `json:"name"`
	DisplayName             string         `json:"display_name"`
	RepoName                string         `json:"repository_name"`
	RepoBranch              string         `json:"repository_branch"`
	Commit                  string         `json:"commit"`
	Description             string         `json:"description"`
	Url                     string         `json:"url"`
	ModifiedBy              string         `json:"modified_by"`
	LastModified            string         `json:"last_modified"`
	Published               bool           `json:"enabled"`
	Inputs                  []Input        `json:"inputs"`
	Tags                    []BlueprintTag `json:"tags"`
	Policies                Policies       `json:"policies"`
	NumOfActiveEnvironments int32          `json:"num_of_active_environments"`
}

type Input struct {
	Name           string   `json:"name"`
	PossibleValues []string `json:"possible_values"`
	DefaultValue   string   `json:"default_value"`
	Description    string   `json:"description"`
}

type BlueprintTag struct {
	Name           string   `json:"name"`
	DefaultValue   string   `json:"default_value"`
	PossibleValues []string `json:"possible_values"`
	Description    string   `json:"description"`
}

type Policies struct {
	MaxDuration           string `json:"max_duration"`
	DefaultDuration       string `json:"default_duration"`
	DefaultExtend         string `json:"default_extend"`
	MaxActiveEnvironments int32  `json:"max_active_environments"`
	AlwaysOn              bool   `json:"always_on"`
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

type GitlabEnterpriseRepoSpaceAssociation struct {
	BaseUrl string `json:"base_url"`
	Name    string `json:"repository_name"`
	URL     string `json:"repository_url"`
	Token   string `json:"token"`
	Branch  string `json:"branch"`
}

type CodeCommitRepoSpaceAssociation struct {
	URL        string `json:"repository_url"`
	RoleArn    string `json:"role_arn"`
	Region     string `json:"aws_region"`
	Branch     string `json:"branch"`
	Name       string `json:"repository_name"`
	ExternalId string `json:"external_id"`
	Username   string `json:"git_username"`
	Password   string `json:"git_password"`
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

type Environment struct {
	ReadOnly          bool                         `json:"read_only"`
	IsWorkflow        bool                         `json:"is_workflow"`
	EnvironmentId     string                       `json:"environment_id"`
	LastUsed          string                       `json:"last_used"`
	IsEAC             bool                         `json:"eac_synced"`
	Details           EnvironmentDetails           `json:"details"`
	Owner             EnvironmentOwner             `json:"owner"`
	Initiator         EnvironmentInitiator         `json:"initiator"`
	CollaboratorsInfo EnvironmentCollaboratorsInfo `json:"collaborators_info"`
}

type EnvironmentDetails struct {
	Id             string                `json:"id"`
	ComputedStatus string                `json:"computed_status"`
	Definition     EnvironmentDefinition `json:"definition"`
	State          EnvironmentState      `json:"state"`
}

type EnvironmentDefinition struct {
	Metadata EnvironmentMetadata `json:"metadata"`
	Inputs   []KeyValuePair      `json:"inputs"`
	Tags     []KeyValuePair      `json:"tags"`
}

type EnvironmentState struct {
	Outputs   []KeyValuePair `json:"outputs"`
	IsEac     bool           `json:"eac_synced"`
	Execution Execution      `json:"execution"`
	Errors    []Error        `json:"errors"`
	Grains    []Grain        `json:"grains"`
}

type KeyValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Error struct {
	Message string `json:"message"`
}

type Execution struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type EnvironmentMetadata struct {
	Name                    string `json:"name"`
	BlueprintName           string `json:"blueprint_name"`
	BlueprintCommit         string `json:"blueprint_commit"`
	BlueprintRepositoryName string `json:"repository_name"`
	SpaceName               string `json:"space_name"`
}

type EnvironmentOwner struct {
	OwnerEmail string `json:"email"`
}

type EnvironmentInitiator struct {
	InitiatorEmail string `json:"email"`
}

type EnvironmentCollaboratorsInfo struct {
	Collaborators   []EnvironmentCollaborator `json:"collaborators"`
	AllSpaceMembers bool                      `json:"all_space_members"`
}

type EnvironmentCollaborator struct {
	Email string `json:"email"`
}

type Grain struct {
	Name    string        `json:"name"`
	Kind    string        `json:"kind"`
	Id      string        `json:"id"`
	Path    string        `json:"path"`
	State   GrainState    `json:"state"`
	Sources []GrainSource `json:"sources"`
}

type GrainState struct {
	CurrentState string `json:"current_state"`
}

type GrainSource struct {
	Store        string `json:"store"`
	Path         string `json:"path"`
	Branch       string `json:"branch"`
	Commit       string `json:"commit"`
	IsLastCommit bool   `json:"is_last_commit"`
}

type IntrospectionItem struct {
	GrainPath        string         `json:"grain_path"`
	GrainType        string         `json:"grain_type"`
	ResourceName     string         `json:"resource_name"`
	ResourceType     string         `json:"resource_type"`
	ResourceCategory string         `json:"resource_category"`
	Status           string         `json:"status"`
	Alias            string         `json:"alias"`
	HasRunningAction bool           `json:"has_running_action"`
	Attributes       []KeyValuePair `json:"attributes"`
	CustomIcon       string         `json:"custom_icon"`
}

type Workflow struct {
	Yaml            string          `json:"yaml"`
	DisplayName     string          `json:"display_name"`
	Description     string          `json:"description"`
	SpaceDefinition spaceDefinition `json:"space_definition"`
	Name            string          `json:"name"`
}

type spaceDefinition struct {
	EnforcedOnAllSpaces bool     `json:"enforced_on_all_spaces"`
	SpecificSpaces      []string `json:"specific_spaces"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CatalogItemLabelsRequest struct {
	BlueprintName  string   `json:"blueprint_name"`
	RepositoryName string   `json:"repository_name"`
	Labels         []string `json:"labels"`
}

type SpaceParameterRequest struct {
	Value       string `json:"value"`
	Sensitive   bool   `json:"sensitive"`
	Description string `json:"description"`
}
