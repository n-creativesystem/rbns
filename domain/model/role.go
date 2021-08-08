package model

type OrganizationUserRole struct {
	organization Organization
	userKey      requiredString
	roleId       id
}

func (m *OrganizationUserRole) GetOrganization() Organization {
	return m.organization
}

func (m *OrganizationUserRole) GetUserKey() string {
	return *m.userKey.Value()
}

func (m *OrganizationUserRole) GetRoleId() string {
	return *m.roleId.Value()
}

type OrganizationUserRoles []OrganizationUserRole

func (arr OrganizationUserRoles) Copy() OrganizationUserRoles {
	cArr := make([]OrganizationUserRole, len(arr))
	copy(cArr, arr)
	return cArr
}

func NewOrganizationUserRole(organization Organization, userKey, roleId string) (*OrganizationUserRole, error) {
	key, err := newKey(userKey)
	if err != nil {
		return nil, err
	}
	roleID, err := newID(roleId)
	if err != nil {
		return nil, err
	}
	return &OrganizationUserRole{
		organization: organization,
		userKey:      *key,
		roleId:       *roleID,
	}, nil
}

type Role struct {
	model
	description           string
	permissions           Permissions
	organizationUserRoles OrganizationUserRoles
}

func (r *Role) GetDescription() string {
	return r.description
}

func (r *Role) GetPermissions() Permissions {
	return r.permissions.Copy()
}

func (r *Role) GetOrganizationUserRoles() OrganizationUserRoles {
	return r.organizationUserRoles.Copy()
}

// AddPermission is only the permission that exists in the database
func (r *Role) AddPermission(permissions ...Permission) {
	nowPermissions := r.permissions.Copy()
	for _, nowPermission := range nowPermissions {
		for _, permission := range permissions {
			if !nowPermission.equals(permission.model) {
				r.permissions = append(r.permissions, permission)
			}
		}
	}
}

func NewRole(id, name, description string, permissions Permissions, organizationUserRoles ...OrganizationUserRole) (*Role, error) {
	vId, err := newID(id)
	if err != nil {
		return nil, err
	}
	vName, err := newName(name)
	if err != nil {
		return nil, err
	}
	return &Role{
		model: model{
			id:   *vId,
			name: *vName,
		},
		description:           description,
		permissions:           permissions.Copy(),
		organizationUserRoles: OrganizationUserRoles(organizationUserRoles).Copy(),
	}, nil
}

type Roles []Role

func (arr Roles) Copy() Roles {
	cArr := make([]Role, len(arr))
	copy(cArr, arr)
	return cArr
}
