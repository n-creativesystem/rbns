package model

type User struct {
	ID   string
	Name string

	Organizations []Organization
	Permissions   []Permission
	Roles         []Role
}

type GetUserQuery struct {
	Result []User
}

type GetUserByIDQuery struct {
	PrimaryQuery

	Result *User
}

type GetUserByIDsQuery struct {
	Query []PrimaryQuery

	Result []User
}

type AddUserCommand struct {
	PrimaryCommand
	Name string

	Result *User
}

type DeleteUserCommand struct {
	PrimaryCommand
}

type AddUserRoleCommand struct {
	Organization *Organization

	PrimaryCommand
	// Roles 存在するPermissionIDをセットする
	Roles []Role
}

type DeleteUserRoleCommand struct {
	Organization *Organization

	PrimaryCommand
	// Roles 存在するPermissionIDをセットする
	Roles []Role
}
