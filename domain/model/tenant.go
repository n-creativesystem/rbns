package model

type Tenant struct {
	ID   string
	Name string
}

type AddTenantCommand struct {
	Name string

	Result *Tenant
}

type DeleteTenantCommand struct {
	ID string
}

type GetTenantQuery struct {
	Result []*Tenant
}

type GetTenantByNameQuery struct {
	Name string

	Result *Tenant
}

type GetTenantByIdQuery struct {
	ID string

	Result *Tenant
}

type AddTenantAndLoginUserCommand struct {
	Tenant    *Tenant
	LoginUser *LoginUser
}
