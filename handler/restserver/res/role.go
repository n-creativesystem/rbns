package res

import (
	"github.com/n-creativesystem/rbns/domain/model"
)

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions
	OrganizationUsers []OrganizationUser `json:"organizationUsers"`
}

type Roles struct {
	Roles []Role `json:"roles"`
}

func NewRole(role model.Role) Role {
	mOrganizationRoles := role.GetOrganizationUserRoles()
	userKeys := make([]OrganizationUser, len(mOrganizationRoles))
	for idx, orgUserRole := range mOrganizationRoles {
		org := orgUserRole.GetOrganization()
		userKeys[idx] = OrganizationUser{
			UserKey:                 orgUserRole.GetUserKey(),
			OrganizationId:          *org.GetID(),
			OrganizationName:        *org.GetName(),
			OrganizationDescription: org.GetDescription(),
		}
	}

	return Role{
		Id:                *role.GetID(),
		Name:              *role.GetName(),
		Description:       role.GetDescription(),
		Permissions:       NewPermissions(role.GetPermissions()),
		OrganizationUsers: userKeys,
	}
}

func NewRoles(roles model.Roles) Roles {
	rs := make([]Role, len(roles))
	for idx, r := range roles {
		rs[idx] = NewRole(r)
	}
	return Roles{
		Roles: rs,
	}
}
