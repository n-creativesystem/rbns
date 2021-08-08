package model

type Permission struct {
	model
	description string
}

func (p *Permission) GetDescription() string {
	return p.description
}

func NewPermission(id, name, description string) (*Permission, error) {
	vId, err := newID(id)
	if err != nil {
		return nil, err
	}
	vName, err := newName(name)
	if err != nil {
		return nil, err
	}
	return &Permission{
		model: model{
			id:   *vId,
			name: *vName,
		},
		description: description,
	}, nil
}

type Permissions []Permission

func (arr Permissions) Copy() Permissions {
	cArr := make([]Permission, len(arr))
	copy(cArr, arr)
	return cArr
}

type ResourceCheck struct {
	result  bool
	message string
}

func NewResourceCheck(result bool, message string) *ResourceCheck {
	return &ResourceCheck{
		result:  result,
		message: message,
	}
}

func (r *ResourceCheck) GetMsg() string {
	return r.message
}

func (r *ResourceCheck) IsOk() bool {
	return r.result
}
