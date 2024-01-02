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
