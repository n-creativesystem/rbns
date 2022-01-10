package protoconv

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/protobuf"
)

func NewRoleEntityByModel(role model.Role) *protobuf.RoleEntity {
	mPermissions := role.Permissions
	permissions := make([]*protobuf.PermissionEntity, len(mPermissions))
	for idx, permission := range mPermissions {
		permissions[idx] = NewPermissionEntityByModel(permission)
	}
	// mOrganizationRoles := role.Organizations
	// userKeys := make([]*protobuf.OrganizationUser, len(mOrganizationRoles))
	// for idx, orgUserRole := range mOrganizationRoles {
	// 	org := orgUserRole.GetOrganization()
	// 	userKeys[idx] = &protobuf.OrganizationUser{
	// 		UserKey:                 orgUserRole.GetUserKey(),
	// 		OrganizationId:          org.GetID().String(),
	// 		OrganizationName:        *org.GetName(),
	// 		OrganizationDescription: org.GetDescription(),
	// 	}
	// }
	return &protobuf.RoleEntity{
		Id:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
		// OrganizationUsers: userKeys,
	}
}

func NewPermissionEntityByModel(permission model.Permission) *protobuf.PermissionEntity {
	return &protobuf.PermissionEntity{
		Id:          permission.ID.String(),
		Name:        permission.Name,
		Description: permission.Description,
	}
}

func NewOrganizationEntityByModel(organization model.Organization) *protobuf.OrganizationEntity {
	mUsers := organization.Users
	users := make([]*protobuf.UserEntity, len(mUsers))
	for idx, user := range mUsers {
		users[idx] = NewUserEntityByModel(user)
	}
	return &protobuf.OrganizationEntity{
		Id:          organization.ID.String(),
		Name:        organization.Name,
		Description: organization.Description,
		Users:       users,
	}
}

func NewUserEntityByModel(user model.User) *protobuf.UserEntity {
	mRoles := user.Roles
	roles := make([]*protobuf.RoleEntity, len(mRoles))
	for idx, role := range mRoles {
		roles[idx] = NewRoleEntityByModel(role)
	}
	mPermissions := user.Permissions
	permissions := make([]*protobuf.PermissionEntity, len(mPermissions))
	for idx, permission := range mPermissions {
		permissions[idx] = NewPermissionEntityByModel(permission)
	}
	return &protobuf.UserEntity{
		Id:          user.ID,
		Roles:       roles,
		Permissions: permissions,
	}
}

// func NewResourceByModel(resource model.Resource) *protobuf.ResourceResponse {
// 	result := &protobuf.ResourceResponse{
// 		Id:          resource.ID,
// 		Description: resource.Description,
// 		Permissions: make([]*protobuf.PermissionEntity, 0, len(resource.Permissions)),
// 	}
// 	permissions := resource.Permissions.Copy()
// 	for _, permission := range permissions {
// 		p := NewPermissionEntityByModel(permission)
// 		result.Permissions = append(result.Permissions, p)
// 	}
// 	return result
// }
