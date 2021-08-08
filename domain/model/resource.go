package model

type Resource struct {
	node Node
}

func NewResourceByMethod(method string, uri, permissions []string) *Resource {
	node := Node{}
	for _, u := range uri {
		node.Add(method, u, permissions...)
	}
	return &Resource{
		node: node,
	}
}

func NewResource(method, uri string, permissions ...string) *Resource {
	return NewResourceByMethod(method, []string{uri}, permissions)
}

func (r *Resource) Check(method, uri string, user *User) bool {
	permissions := r.node.Get(method, uri)
	for _, permission := range permissions {
		for _, up := range user.permissions {
			if *up.name.Value() == permission {
				return true
			}
		}
	}
	return false
}
