package protoconv

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/protobuf"
)

func NewRoleEntityByModel(role model.Role) *protobuf.RoleEntity {
	mPermissions := role.GetPermissions()
	permissions := make([]*protobuf.PermissionEntity, len(mPermissions))
	for idx, permission := range mPermissions {
		permissions[idx] = NewPermissionEntityByModel(permission)
	}
	mOrganizationRoles := role.GetOrganizationUserRoles()
	userKeys := make([]*protobuf.OrganizationUser, len(mOrganizationRoles))
	for idx, orgUserRole := range mOrganizationRoles {
		org := orgUserRole.GetOrganization()
		userKeys[idx] = &protobuf.OrganizationUser{
			UserKey:                 orgUserRole.GetUserKey(),
			OrganizationId:          *org.GetID(),
			OrganizationName:        *org.GetName(),
			OrganizationDescription: org.GetDescription(),
		}
	}
	return &protobuf.RoleEntity{
		Id:                *role.GetID(),
		Name:              *role.GetName(),
		Description:       role.GetDescription(),
		Permissions:       permissions,
		OrganizationUsers: userKeys,
	}
}

func NewPermissionEntityByModel(permission model.Permission) *protobuf.PermissionEntity {
	return &protobuf.PermissionEntity{
		Id:          *permission.GetID(),
		Name:        *permission.GetName(),
		Description: permission.GetDescription(),
	}
}

func NewOrganizationEntityByModel(organization model.Organization) *protobuf.OrganizationEntity {
	mUsers := organization.GetUsers()
	users := make([]*protobuf.UserEntity, len(mUsers))
	for idx, user := range mUsers {
		users[idx] = NewUserEntityByModel(user)
	}
	return &protobuf.OrganizationEntity{
		Id:          *organization.GetID(),
		Name:        *organization.GetName(),
		Description: organization.GetDescription(),
		Users:       users,
	}
}

func NewUserEntityByModel(user model.User) *protobuf.UserEntity {
	mRoles := user.GetRole()
	roles := make([]*protobuf.RoleEntity, len(mRoles))
	for idx, role := range mRoles {
		roles[idx] = NewRoleEntityByModel(role)
	}
	mPermissions := user.GetPermission()
	permissions := make([]*protobuf.PermissionEntity, len(mPermissions))
	for idx, permission := range mPermissions {
		permissions[idx] = NewPermissionEntityByModel(permission)
	}
	return &protobuf.UserEntity{
		Key:         user.GetKey(),
		Roles:       roles,
		Permissions: permissions,
	}
}

func NewResourceByModel(resource model.Resource) *protobuf.ResourceResponse {
	result := &protobuf.ResourceResponse{
		Id:          resource.ID,
		Description: resource.Description,
		Permissions: make([]*protobuf.PermissionEntity, 0, len(resource.Permissions)),
	}
	permissions := resource.Permissions.Copy()
	for _, permission := range permissions {
		p := NewPermissionEntityByModel(permission)
		result.Permissions = append(result.Permissions, p)
	}
	return result
}
