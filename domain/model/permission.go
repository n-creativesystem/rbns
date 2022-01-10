package model

type Permission struct {
	ID          ID
	Name        string
	Description string

	Roles []Role
}

type GetPermissionQuery struct {
	Result []Permission
}

type GetPermissionByIDQuery struct {
	PrimaryCommand

	Result *Permission
}

type GetPermissionByIDsQuery struct {
	Query []PrimaryCommand

	Result []Permission
}

type GetPermissionByNameQuery struct {
	Name Name

	Result *Permission
}

type CountPermissionByNameQuery struct {
	CountNameQuery
}

type AddPermissionCommand struct {
	Name        Name
	Description string

	Result *Permission
}

type AddPermissionCommands struct {
	AddPermissions []AddPermissionCommand
}

type UpdatePermissionCommand struct {
	PrimaryCommand
	Name        Name
	Description string

	Result *Permission
}

type DeletePermissionCommand struct {
	PrimaryCommand
}
