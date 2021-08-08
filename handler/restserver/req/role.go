package req

type RoleCreateBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type RolesCreateBody struct {
	Roles []RoleCreateBody `json:"roles"`
}

type RoleUpdateBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionKey struct {
	Id string `json:"id"`
}
type RolePermission struct {
	Permissions []PermissionKey `json:"permissions"`
}
