package req

type RoleBodyes struct {
	Id string `json:"id"`
}
type UserCreateBody struct {
	Key            string       `json:"key"`
	OrganizationId string       `json:"organization_id"`
	Roles          []RoleBodyes `json:"roles"`
}

type UserUpdateBody struct {
	Roles []RoleBodyes `json:"roles"`
}
