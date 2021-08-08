package protoconv

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/proto"
)

func NewRoleEntityByModel(role model.Role) *proto.RoleEntity {
	mPermissions := role.GetPermissions()
	permissions := make([]*proto.PermissionEntity, len(mPermissions))
	for idx, permission := range mPermissions {
		permissions[idx] = NewPermissionEntityByModel(permission)
	}
	mOrganizationRoles := role.GetOrganizationUserRoles()
	userKeys := make([]*proto.OrganizationUser, len(mOrganizationRoles))
	for idx, orgUserRole := range mOrganizationRoles {
		org := orgUserRole.GetOrganization()
		userKeys[idx] = &proto.OrganizationUser{
			UserKey:                 orgUserRole.GetUserKey(),
			OrganizationId:          *org.GetID(),
			OrganizationName:        *org.GetName(),
			OrganizationDescription: org.GetDescription(),
		}
	}
	return &proto.RoleEntity{
		Id:                *role.GetID(),
		Name:              *role.GetName(),
		Description:       role.GetDescription(),
		Permissions:       permissions,
		OrganizationUsers: userKeys,
	}
}

func NewPermissionEntityByModel(permission model.Permission) *proto.PermissionEntity {
	return &proto.PermissionEntity{
		Id:          *permission.GetID(),
		Name:        *permission.GetName(),
		Description: permission.GetDescription(),
	}
}

func NewOrganizationEntityByModel(organization model.Organization) *proto.OrganizationEntity {
	mUsers := organization.GetUsers()
	users := make([]*proto.UserEntity, len(mUsers))
	for idx, user := range mUsers {
		users[idx] = NewUserEntityByModel(user)
	}
	return &proto.OrganizationEntity{
		Id:          *organization.GetID(),
		Name:        *organization.GetName(),
		Description: organization.GetDescription(),
		Users:       users,
	}
}

func NewUserEntityByModel(user model.User) *proto.UserEntity {
	mRoles := user.GetRole()
	roles := make([]*proto.RoleEntity, len(mRoles))
	for idx, role := range mRoles {
		roles[idx] = NewRoleEntityByModel(role)
	}
	mPermissions := user.GetPermission()
	permissions := make([]*proto.PermissionEntity, len(mPermissions))
	for idx, permission := range mPermissions {
		permissions[idx] = NewPermissionEntityByModel(permission)
	}
	return &proto.UserEntity{
		Key:         user.GetKey(),
		Roles:       roles,
		Permissions: permissions,
	}
}
