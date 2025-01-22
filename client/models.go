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
	MaxActiveEnvironments *int32 `json:"max_active_environments"`
	AlwaysOn              bool   `json:"always_on"`
	AllowScheduling       bool   `json:"allow_scheduling"`
}

type UserSpaceAssociation struct {
	Email     string `json:"email"`
	SpaceRole string `json:"space_role"`
}

type Space struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	NumOfUsers  int32  `json:"users_count"`
	NumOfGroups int32  `json:"groups_count"`
}

type AgentSpaceAssociation struct {
	Type                  string `json:"type"`
	DefaultNamespace      string `json:"namespace"`
	DefaultServiceAccount string `json:"service_account"`
}

type RepoSpaceAssociation struct {
	URL         string  `json:"repository_url"`
	AccessToken *string `json:"access_token"`
	Type        string  `json:"repository_type"`
	Branch      string  `json:"branch"`
	Name        string  `json:"repository_name"`
}

type RepoSpaceAssociationWithCredentials struct {
	URL            string  `json:"repository_url"`
	Token          *string `json:"token"`
	Type           string  `json:"repository_type"`
	Branch         string  `json:"branch"`
	Name           string  `json:"repository_name"`
	CredentialName *string `json:"credential_name"`
}

type GitlabEnterpriseRepoSpaceAssociation struct {
	Name           string  `json:"repository_name"`
	URL            string  `json:"repository_url"`
	Token          *string `json:"token"`
	Branch         string  `json:"branch"`
	CredentialName string  `json:"credential_name"`
}

type CodeCommitRepoSpaceAssociation struct {
	URL            string `json:"repository_url"`
	RoleArn        string `json:"role_arn"`
	Region         string `json:"region"`
	Branch         string `json:"branch"`
	Name           string `json:"repository_name"`
	ExternalId     string `json:"external_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	CredentialName string `json:"credential_name"`
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Tag struct {
	Name           string   `json:"name"`
	Value          string   `json:"value"`
	Scope          string   `json:"scope"`
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

type AwsCostTarget struct {
	Name       string `json:"name"`
	NewName    string `json:"new_name"`
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
	Type        string  `json:"type"`
	Description string  `json:"description"`
	WebHook     *string `json:"web_hook,omitempty"`
	Token       *string `json:"token,omitempty"`
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
	Inputs   []NameValuePair     `json:"inputs"`
	Tags     []NameValuePair     `json:"tags"`
	Labels   []KeyValuePair      `json:"labels"`
}

type EnvironmentState struct {
	Outputs      []NameValuePair `json:"outputs"`
	IsEac        bool            `json:"eac_synced"`
	Execution    Execution       `json:"execution"`
	Errors       []Error         `json:"errors"`
	Grains       []Grain         `json:"grains"`
	CurrentState string          `json:"current_state"`
}

type NameValuePair struct {
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

// Data Source Collaborators.
type EnvironmentCollaboratorsInfo struct {
	Collaborators   []EnvironmentCollaborator `json:"collaborators"`
	AllSpaceMembers bool                      `json:"all_space_members"`
}

// Resource Collaborators.
type Collaborators struct {
	Collaborators   []string `json:"collaborators_emails"`
	AllSpaceMembers bool     `json:"all_space_members"`
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
	GrainPath        string          `json:"grain_path"`
	GrainType        string          `json:"grain_type"`
	ResourceName     string          `json:"resource_name"`
	ResourceType     string          `json:"resource_type"`
	ResourceCategory string          `json:"resource_category"`
	Status           string          `json:"status"`
	Alias            string          `json:"alias"`
	HasRunningAction bool            `json:"has_running_action"`
	Attributes       []NameValuePair `json:"attributes"`
	CustomIcon       string          `json:"custom_icon"`
}

type Workflow struct {
	Yaml            string          `json:"yaml"`
	DisplayName     string          `json:"display_name"`
	Description     string          `json:"description"`
	SpaceDefinition spaceDefinition `json:"space_definition"`
	Name            string          `json:"name"`
}

type EnvironmentWorkflow struct {
	Name            string            `json:"name"`
	Schedules       []Schedule        `json:"schedules"`
	Reminder        int64             `json:"reminder"`
	InputsOverrides map[string]string `json:"inputs_overrides"`
}

type Schedule struct {
	Scheduler  string `json:"scheduler"`
	Overridden bool   `json:"overridden"`
}

type spaceDefinition struct {
	EnforcedOnAllSpaces bool     `json:"enforced_on_all_spaces"`
	SpecificSpaces      []string `json:"specific_spaces"`
}

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	QuickFilter bool   `json:"quick_filter"`
}

type LabelRequest struct {
	OriginalName string `json:"original_name"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	QuickFilter  bool   `json:"quick_filter"`
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

type EnvironmentLabelsUpdateRequest struct {
	EnvironmentId string         `json:"environment_id"`
	SpaceName     string         `json:"space_name"`
	AddedLabels   []KeyValuePair `json:"added_labels"`
	RemovedLabels []KeyValuePair `json:"removed_labels"`
}

type SpaceCredentials struct {
	SpaceName       string         `json:"space_name"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	CredentialData  CredentialData `json:"credential_data"`
	CloudType       string         `json:"cloud_type"`
	CloudIdentifier string         `json:"cloud_identifier"`
}

type AccountCredentials struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	CredentialData    CredentialData `json:"credential_data"`
	CloudType         string         `json:"cloud_type"`
	CloudIdentifier   string         `json:"cloud_identifier"`
	AllowedSpaceNames []string       `json:"allowed_space_names,omitempty"`
	AllSpacesAllowed  bool           `json:"all_spaces_allowed"`
}

type CredentialData struct {
	Token  *string `json:"token,omitempty"`
	Key    *string `json:"key,omitempty"`
	Secret *string `json:"secret,omitempty"`
	Type   string  `json:"type"`
}

type BlueprintSource struct {
	BlueprintName  *string `tfsdk:"blueprint_name"`
	RepositoryName *string `tfsdk:"repository_name"`
	Branch         *string `tfsdk:"branch"`
	Commit         *string `tfsdk:"commit"`
}

type EnvironmentRequest struct {
	EnvironmentName  string                `json:"environment_name"`
	BlueprintName    string                `json:"blueprint_name"`
	OwnerEmail       string                `json:"owner_email"`
	Description      string                `json:"description"`
	Inputs           map[string]string     `json:"inputs"`
	Tags             map[string]string     `json:"tags"`
	Collaborators    Collaborators         `json:"collaborators"`
	Automation       bool                  `json:"automation"`
	ScheduledEndTime string                `json:"scheduled_end_time"`
	Duration         string                `json:"duration"`
	Id               string                `json:"id"`
	BlueprintSource  BlueprintSource       `json:"blueprint_source"`
	Workflows        []EnvironmentWorkflow `json:"workflows"`
}

type WorkflowRequest struct {
	BlueprintName  string `json:"blueprint_name"`
	RepositoryName string `json:"repository_name"`
	SpaceName      string `json:"space_name"`
	LaunchAllowed  bool   `json:"launch_allowed"`
}

type SpaceWorkflow struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}

type BlueprintDisplayNameRequest struct {
	BlueprintName  string `json:"blueprint_name"`
	RepositoryName string `json:"repository_name"`
	DisplayName    string `json:"display_name"`
}

type TorqueSpaceCustomIcon struct {
	FileName string `json:"file_name"`
	Key      string `json:"key"`
}

type TorqueInputSource struct {
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Details       InputSourceDetails `json:"details"`
	AllowedSpaces AllowedSpaces      `json:"allowed_spaces"`
}

type InputSourceDetails struct {
	CredentialName     string            `json:"credential_name"`
	BucketName         *OverridableValue `json:"bucket_name,omitempty"`
	StorageAccountName *OverridableValue `json:"storage_account_name,omitempty"`
	ContainerName      *OverridableValue `json:"container_name,omitempty"`
	BlobName           *OverridableValue `json:"blob_name,omitempty"`
	PathPrefix         *OverridableValue `json:"path_prefix"`
	FilterPattern      OverridableValue  `json:"filter_pattern"`
	ObjectKey          *OverridableValue `json:"object_key"`
	ContentFormat      *ContentFormat    `json:"content_format"`
	Type               string            `json:"type"`
}

type AllowedSpaces struct {
	AllSpaces      bool     `json:"all_spaces"`
	SpecificSpaces []string `json:"specific_spaces"`
}

type OverridableValue struct {
	Overridable bool   `json:"overridable"`
	Value       string `json:"value"`
}

type ContentFormat struct {
	DisplayJsonPath OverridableValue `json:"display_json_path"`
	JsonPath        OverridableValue `json:"json_path"`
	Type            string           `json:"type"`
}

type ResourceInventory struct {
	Credentials string                   `json:"credentials"`
	Details     ResourceInventoryDetails `json:"details"`
}

type ResourceInventoryDetails struct {
	Type    string  `json:"type"`
	ViewArn *string `json:"view_arn"`
}

type DeploymentEngine struct {
	Name                   string        `json:"name"`
	Description            string        `json:"description"`
	Type                   string        `json:"type"`
	AuthToken              string        `json:"auth_token"`
	AgentName              string        `json:"agent_name"`
	ServerUrl              string        `json:"server_url"`
	PollingIntervalSeconds int32         `json:"polling_interval_seconds"`
	AllowedSpaces          AllowedSpaces `json:"allowed_spaces"`
}

type DeploymentEngineRead struct {
	Name                   string          `json:"name"`
	Description            string          `json:"description"`
	Type                   string          `json:"type"`
	AuthToken              string          `json:"auth_token"`
	Agent                  AgentDetails    `json:"agent"`
	ServerUrl              string          `json:"server_url"`
	PollingIntervalSeconds PollingInterval `json:"polling_interval_seconds"`
	AllowedSpaces          AllowedSpaces   `json:"allowed_spaces"`
}

type PollingInterval struct {
	Default float32 `json:"default"`
	Value   float32 `json:"value"`
}

type AgentDetails struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ApprovalChannel struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Details     ApprovalChannelDetails `json:"details"`
}

type ApprovalChannelDetails struct {
	Type           string     `json:"type"`
	Approvers      []Approver `json:"approvers"`
	Approver       *Approver  `json:"approver,omitempty"`
	Headers        *string    `json:"headers,omitempty"`
	BaseUrl        *string    `json:"base_url,omitempty"`
	UserName       *string    `json:"user_name,omitempty"`
	Password       *string    `json:"password,omitempty"`
	WebhookAddress *string    `json:"webhook_address,omitempty"`
}

type Approver struct {
	UserEmail string `json:"user_email"`
}

type Audit struct {
	Type       string           `json:"type"`
	Properties *AuditProperties `json:"properties,omitempty"`
}

type AuditProperties struct {
	Url         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Certificate string `json:"certificate"`
}
