package model

type User struct {
	key         requiredString
	permissions Permissions
	roles       Roles
}

func NewUser(key string, roles Roles, permissions Permissions) (*User, error) {
	v, err := newKey(key)
	if err != nil {
		return nil, err
	}
	return &User{
		key:         *v,
		permissions: permissions.Copy(),
		roles:       roles.Copy(),
	}, nil
}

func (u *User) GetKey() string {
	return string(u.key)
}

func (u *User) GetPermission() []Permission {
	return u.permissions.Copy()
}

func (u *User) GetRole() []Role {
	return u.roles.Copy()
}

func (u *User) AddRole(roles ...Role) {
	nowRoles := u.GetRole()
	for _, role := range roles {
		isExists := false
		for _, nowRole := range nowRoles {
			isExists = nowRole.equals(role.model)
		}
		if !isExists {
			u.roles = append(u.roles, role)
		}
	}
}

func (u *User) IsContainsPermission(id ID) bool {
	permissions := u.GetPermission()
	for _, permission := range permissions {
		if permission.id.equals(id) {
			return true
		}
	}
	return false
}

func (u *User) IsContainsPermissionByName(name Name) bool {
	permissions := u.GetPermission()
	for _, permission := range permissions {
		if permission.name.equals(name) {
			return true
		}
	}
	return false
}

type Users []User

func (arr Users) Copy() Users {
	cArr := make([]User, len(arr))
	copy(cArr, arr)
	return cArr
}
