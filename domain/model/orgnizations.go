package model

type Organization struct {
	ID          ID
	Name        string
	Description string
	Users       []User
	Roles       []Role
}

func (org *Organization) AddUser(id ID, name Name) *Organization {
	org.Users = append(org.Users, User{
		ID:   id.String(),
		Name: name.String(),
	})
	return org
}

type GetOrganizationQuery struct {
	Result []Organization
}

type GetOrganizationByIDQuery struct {
	PrimaryQuery

	Result *Organization
}

type GetOrganizationByNameQuery struct {
	Name Name

	Result *Organization
}

type CountOrganizationByNameQuery struct {
	CountNameQuery
}

type AddOrganizationCommand struct {
	Name        Name
	Description string

	Result *Organization
}

type UpdateOrganizationCommand struct {
	PrimaryCommand
	Name        Name
	Description string

	Result *Organization
}

type DeleteOrganizationCommand struct {
	PrimaryCommand
}

type AddOrganizationUserCommand struct {
	PrimaryCommand
	User []User
}

type DeleteOrganizationUserCommand struct {
	PrimaryCommand
	User []User
}
