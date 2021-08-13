package req

type PermissionCreateBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type PermissionsCreateBody struct {
	Permissions []PermissionCreateBody `json:"permissions"`
}

type PermissionUpdateBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionsCheckBody struct {
	UserKey          string   `json:"user_key" header:"x-user-key"`
	OrganizationName string   `json:"organization_name" header:"x-organization"`
	PermissionNames  []string `json:"permission_names" header:"x-permissions"`
}
