package model

type Resource struct {
	ID          string
	Description string
	Permissions Permissions
}

func NewResource(id string, permissions ...Permission) *Resource {
	p := make(Permissions, 0, len(permissions))
	for _, permission := range permissions {
		p = append(p, permission)
	}
	return &Resource{
		ID:          id,
		Permissions: p,
	}
}

func (r *Resource) Add(permission Permission) {
	r.Permissions = append(r.Permissions, permission)
}
