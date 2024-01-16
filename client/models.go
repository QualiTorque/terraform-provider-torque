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

type CatalogItemRequest struct {
	BlueprintName  string `json:"blueprint_name"`
	RepositoryName string `json:"repository_name"`
}
