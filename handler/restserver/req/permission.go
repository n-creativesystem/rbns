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
