package model

type Role struct {
	ID          ID
	Name        string
	Description string

	Permissions   []Permission
	Organizations []Organization
}

// AddPermission
func (r *Role) AddPermission(id ID) *Role {
	r.Permissions = append(r.Permissions, Permission{
		ID: id,
	})
	return r
}

// AddOrganization
func (r *Role) AddOrganization(id ID) *Role {
	r.Organizations = append(r.Organizations, Organization{
		ID: id,
	})
	return r
}

type GetRoleQuery struct {
	Organization *Organization
	Result       []Role
}

type GetRoleByIDQuery struct {
	Organization *Organization
	PrimaryQuery

	Result *Role
}

type CountRoleByNameQuery struct {
	CountNameQuery
}

type AddRoleCommand struct {
	Organization *Organization
	Name         Name
	Description  string

	Result *Role
}

type AddRoleCommands struct {
	Organization *Organization
	Roles        []*AddRoleCommand
}

type UpdateRoleCommand struct {
	Organization *Organization
	PrimaryCommand
	Name        Name
	Description string

	Result *Role
}

type DeleteRoleCommand struct {
	Organization *Organization
	PrimaryCommand
}

type DeleteRoleCommands struct {
	Organization *Organization
	Roles        []*DeleteRoleCommand
}

type AddRolePermissionCommand struct {
	*Role
	// Permissions 存在するPermissionsをセットする
	Permissions []Permission
}

type DeleteRolePermissionCommand struct {
	*Role
	// Permissions 存在するPermissionsをセットする
	Permissions []Permission
}
