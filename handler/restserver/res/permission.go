package res

import "github.com/n-creativesystem/rbns/domain/model"

type Permission struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Permissions struct {
	Permissions []Permission `json:"permissions"`
}

func NewPermission(permission model.Permission) Permission {
	return Permission{
		Id:          *permission.GetID(),
		Name:        *permission.GetName(),
		Description: permission.GetDescription(),
	}
}

func NewPermissions(permissions model.Permissions) Permissions {
	ps := make([]Permission, len(permissions))
	for idx, p := range permissions {
		ps[idx] = NewPermission(p)
	}
	return Permissions{
		Permissions: ps,
	}
}

type PermissionCheckRequest struct {
	PermissionName   string `json:"permissionName"`
	UserKey          string `json:"userKey"`
	OrganizationName string `json:"organizationName"`
}
