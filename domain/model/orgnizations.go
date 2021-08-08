package model

type Organization struct {
	model
	description string
	users       Users
}

func (o *Organization) GetDescription() string {
	return o.description
}

func (o *Organization) GetUsers() Users {
	return o.users.Copy()
}

func (o *Organization) IsContainsUsers(userKey Key) (*User, bool) {
	users := o.GetUsers()
	for _, u := range users {
		if u.key.equals(userKey) {
			return &u, true
		}
	}
	return nil, false
}

type Organizations []Organization

func (arr Organizations) Copy() Organizations {
	cArr := make([]Organization, len(arr))
	copy(cArr, arr)
	return cArr
}

func NewOrganization(id, name, description string, users ...User) (*Organization, error) {
	vId, err := newID(id)
	if err != nil {
		return nil, err
	}
	vName, err := newName(name)
	if err != nil {
		return nil, err
	}
	return &Organization{
		model: model{
			id:   *vId,
			name: *vName,
		},
		description: description,
		users:       Users(users).Copy(),
	}, nil
}
